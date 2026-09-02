package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Auth      AuthConfig      `mapstructure:"auth"`
	RateLimit RateLimitConfig `mapstructure:"rate_limiter"`
	Transport TransportConfig `mapstructure:"transport"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Billing   BillingConfig   `mapstructure:"billing"`
	Adaptive  AdaptiveConfig  `mapstructure:"adaptive_transport"`
}

type ServerConfig struct {
	Name            string        `mapstructure:"name"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	GRPCPort        int           `mapstructure:"grpc_port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	EnablePprof     bool          `mapstructure:"enable_pprof"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Name            string        `mapstructure:"name"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type AuthConfig struct {
	AccessSecret     string        `mapstructure:"jwt_access_secret"`
	RefreshSecret    string        `mapstructure:"jwt_refresh_secret"`
	AccessTokenTTL   time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL  time.Duration `mapstructure:"refresh_token_ttl"`
	JWTIssuer        string        `mapstructure:"jwt_issuer"`
	Argon2Time       uint32        `mapstructure:"argon2_time"`
	Argon2Memory     uint32        `mapstructure:"argon2_memory"`
	Argon2Threads    uint8         `mapstructure:"argon2_threads"`
	Argon2KeyLength  uint32        `mapstructure:"argon2_key_length"`
}

type RateLimitConfig struct {
	Enabled          bool          `mapstructure:"enabled"`
	RequestsPerSec   int           `mapstructure:"requests_per_second"`
	Burst            int           `mapstructure:"burst"`
	BanDuration      time.Duration `mapstructure:"ban_duration"`
	MaxLoginAttempts int           `mapstructure:"max_login_attempts"`
	LoginBanDuration time.Duration `mapstructure:"login_ban_duration"`
}

type TransportConfig struct {
	WireGuard WireGuardConfig   `mapstructure:"wireguard"`
	VLESS     VLESSConfig       `mapstructure:"vless_reality"`
	Hysteria2 Hysteria2Config   `mapstructure:"hysteria2"`
	TUIC      TUICConfig        `mapstructure:"tuic"`
}

type WireGuardConfig struct {
	PortRange          string `mapstructure:"port_range"`
	Interface          string `mapstructure:"interface"`
	MTU                int    `mapstructure:"mtu"`
	PersistentKeepalive int   `mapstructure:"persistent_keepalive"`
}

type VLESSConfig struct {
	Port        int      `mapstructure:"port"`
	ServerNames []string `mapstructure:"server_names"`
}

type Hysteria2Config struct {
	Port        int    `mapstructure:"port"`
	UpMbps      int    `mapstructure:"up_mbps"`
	DownMbps    int    `mapstructure:"down_mbps"`
	ObfsPassword string `mapstructure:"obfs_password"`
}

type TUICConfig struct {
	Port int `mapstructure:"port"`
}

type AdaptiveConfig struct {
	Enabled            bool          `mapstructure:"enabled"`
	CheckInterval      time.Duration `mapstructure:"check_interval"`
	MinStableDuration  time.Duration `mapstructure:"min_stable_duration"`
	FallbackCooldown   time.Duration `mapstructure:"fallback_cooldown"`
}

type LoggingConfig struct {
	Level  string   `mapstructure:"level"`
	Format string   `mapstructure:"format"`
	Outputs []string `mapstructure:"outputs"`
}

type BillingConfig struct {
	Stripe StripeConfig `mapstructure:"stripe"`
	Apple  AppleConfig  `mapstructure:"apple"`
	Google GoogleConfig `mapstructure:"google"`
}

type StripeConfig struct {
	SecretKey      string `mapstructure:"secret_key"`
	WebhookSecret  string `mapstructure:"webhook_secret"`
}

type AppleConfig struct {
	SharedSecret string `mapstructure:"shared_secret"`
}

type GoogleConfig struct {
	ServiceAccount string `mapstructure:"service_account"`
	PackageName    string `mapstructure:"package_name"`
}

func Load(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.SetEnvPrefix("NEBULA")

	// Defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.grpc_port", 9090)
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.shutdown_timeout", "30s")

	v.SetDefault("database.port", 5432)
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "5m")

	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 50)
	v.SetDefault("redis.min_idle_conns", 10)

	v.SetDefault("auth.access_token_ttl", "15m")
	v.SetDefault("auth.refresh_token_ttl", "168h")
	v.SetDefault("auth.jwt_issuer", "nebula-vpn")
	v.SetDefault("auth.argon2_time", 3)
	v.SetDefault("auth.argon2_memory", 65536)
	v.SetDefault("auth.argon2_threads", 4)
	v.SetDefault("auth.argon2_key_length", 32)

	v.SetDefault("rate_limiter.enabled", true)
	v.SetDefault("rate_limiter.requests_per_second", 100)
	v.SetDefault("rate_limiter.burst", 200)
	v.SetDefault("rate_limiter.ban_duration", "24h")
	v.SetDefault("rate_limiter.max_login_attempts", 5)
	v.SetDefault("rate_limiter.login_ban_duration", "15m")

	v.SetDefault("transport.wireguard.port_range", "51820-51830")
	v.SetDefault("transport.wireguard.mtu", 1420)

	v.SetDefault("adaptive_transport.enabled", true)
	v.SetDefault("adaptive_transport.check_interval", "30s")
	v.SetDefault("adaptive_transport.min_stable_duration", "120s")

	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.outputs", []string{"stdout"})

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
