package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client     *redis.Client
	enabled    bool
	rps        int
	burst      int
	banDuration time.Duration
	mu         sync.RWMutex
	localStore map[string]*bucket
}

type bucket struct {
	tokens    int
	lastCheck time.Time
	banUntil  time.Time
}

func NewRateLimiter(client *redis.Client, enabled bool, rps, burst int, banDuration time.Duration) *RateLimiter {
	return &RateLimiter{
		client:      client,
		enabled:     enabled,
		rps:         rps,
		burst:       burst,
		banDuration: banDuration,
		localStore:  make(map[string]*bucket),
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, time.Duration) {
	if !rl.enabled {
		return true, 0
	}

	if rl.client != nil {
		return rl.allowRedis(ctx, key)
	}
	return rl.allowLocal(key)
}

func (rl *RateLimiter) allowRedis(ctx context.Context, key string) (bool, time.Duration) {
	// Redis-based sliding window counter using INCR + EXPIRE
	limit := int64(rl.rps)

	now := time.Now().Unix()
	windowKey := fmt.Sprintf("rate_limit:%s:%d", key, now)

	pipe := rl.client.Pipeline()
	count := pipe.Incr(ctx, windowKey)
	pipe.Expire(ctx, windowKey, time.Second)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return true, 0
	}

	if count.Val() > limit {
		retryAfter := time.Duration(int64(time.Second) / limit)
		return false, retryAfter
	}

	return true, 0
}

func (rl *RateLimiter) allowLocal(key string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, exists := rl.localStore[key]
	if !exists {
		rl.localStore[key] = &bucket{
			tokens:    rl.burst - 1,
			lastCheck: time.Now(),
		}
		return true, 0
	}

	if time.Now().Before(b.banUntil) {
		return false, b.banUntil.Sub(time.Now())
	}

	now := time.Now()
	elapsed := now.Sub(b.lastCheck)
	b.lastCheck = now

	// Refill tokens based on elapsed time
	b.tokens += int(elapsed.Seconds() * float64(rl.rps))
	if b.tokens > rl.burst {
		b.tokens = rl.burst
	}

	if b.tokens > 0 {
		b.tokens--
		return true, 0
	}

	b.banUntil = now.Add(rl.banDuration)
	return false, rl.banDuration
}

func (rl *RateLimiter) LoginRateLimit(ctx context.Context, email string) (bool, error) {
	key := fmt.Sprintf("login_attempts:%s", email)
	if rl.client != nil {
		count, err := rl.client.Incr(ctx, key).Result()
		if err != nil {
			return true, nil
		}
		rl.client.Expire(ctx, key, 15*time.Minute)
		return count <= 5, nil
	}
	return true, nil
}

func (rl *RateLimiter) ResetLoginAttempts(ctx context.Context, email string) {
	if rl.client != nil {
		rl.client.Del(ctx, fmt.Sprintf("login_attempts:%s", email))
	}
}
