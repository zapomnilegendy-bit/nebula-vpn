package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/nebula-vpn/backend/internal/config"
	"go.uber.org/zap"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, cfg *config.RedisConfig, logger *zap.Logger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info("connected to Redis", zap.String("addr", cfg.Addr()))

	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

func (r *RedisClient) HealthCheck(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

// Cache helpers for VPN-specific data
func (r *RedisClient) CacheServerList(ctx context.Context, servers interface{}, ttl time.Duration) error {
	return r.Client.Set(ctx, "vpn:servers", servers, ttl).Err()
}

func (r *RedisClient) GetCachedServerList(ctx context.Context) (string, error) {
	return r.Client.Get(ctx, "vpn:servers").Result()
}

func (r *RedisClient) CacheUserSession(ctx context.Context, userID string, token string, ttl time.Duration) error {
	return r.Client.Set(ctx, fmt.Sprintf("session:%s", userID), token, ttl).Err()
}

func (r *RedisClient) InvalidateUserSession(ctx context.Context, userID string) error {
	return r.Client.Del(ctx, fmt.Sprintf("session:%s", userID)).Err()
}

func (r *RedisClient) RecordServerMetric(ctx context.Context, serverID string, field string, value float64) error {
	return r.Client.HIncrByFloat(ctx, fmt.Sprintf("metrics:server:%s", serverID), field, value).Err()
}

func (r *RedisClient) PublishEvent(ctx context.Context, channel string, message interface{}) error {
	return r.Client.Publish(ctx, channel, message).Err()
}
