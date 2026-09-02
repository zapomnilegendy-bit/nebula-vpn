package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nebula-vpn/backend/internal/api/rest"
	"github.com/nebula-vpn/backend/internal/auth"
	"github.com/nebula-vpn/backend/internal/config"
	"github.com/nebula-vpn/backend/internal/server"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/development.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize logger
	var logger *zap.Logger
	if cfg.Logging.Format == "json" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	logger.Info("starting NebulaVPN API server",
		zap.String("host", cfg.Server.Host),
		zap.Int("port", cfg.Server.Port),
	)

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.Auth.AccessSecret,
		cfg.Auth.RefreshSecret,
		cfg.Auth.AccessTokenTTL,
		cfg.Auth.RefreshTokenTTL,
		cfg.Auth.JWTIssuer,
	)

	// Initialize rate limiter
	rateLimiter := auth.NewRateLimiter(
		nil, // Redis client (optional for now)
		cfg.RateLimit.Enabled,
		cfg.RateLimit.RequestsPerSec,
		cfg.RateLimit.Burst,
		cfg.RateLimit.BanDuration,
	)

	// Initialize database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.Database.DSN())
	if err != nil {
		logger.Fatal("failed to parse database config", zap.Error(err))
	}
	poolConfig.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}
	logger.Info("connected to database")

	// Initialize handlers
	handlerCtx := &rest.HandlerContext{
		JWTManager:  jwtManager,
		RateLimiter: rateLimiter,
		Config:      cfg,
		DB:          pool,
	}

	// Create router
	app := rest.NewRouter(handlerCtx)

	// Create and start server
	srv := server.New(cfg, app, logger)
	if err := srv.Start(context.Background()); err != nil {
		logger.Fatal("server error", zap.Error(err))
	}
}
