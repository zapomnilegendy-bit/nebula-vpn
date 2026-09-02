package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nebula-vpn/backend/internal/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	cfg    *config.Config
	app    *fiber.App
	logger *zap.Logger
}

func New(cfg *config.Config, app *fiber.App, logger *zap.Logger) *Server {
	return &Server{
		cfg:    cfg,
		app:    app,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// Start HTTP server
	g.Go(func() error {
		addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
		s.logger.Info("starting HTTP server", zap.String("addr", addr))
		if err := s.app.Listen(addr); err != nil {
			return fmt.Errorf("http server error: %w", err)
		}
		return nil
	})

	// Print stack usage stats periodically
	g.Go(func() error {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				s.logger.Info("memory stats",
					zap.Uint64("alloc_mb", m.Alloc/1024/1024),
					zap.Uint64("total_alloc_mb", m.TotalAlloc/1024/1024),
					zap.Uint64("sys_mb", m.Sys/1024/1024),
					zap.Uint32("num_gc", m.NumGC),
				)
			case <-ctx.Done():
				return nil
			}
		}
	})

	// Graceful shutdown
	g.Go(func() error {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		select {
		case sig := <-sigChan:
			s.logger.Info("shutting down", zap.String("signal", sig.String()))
			cancel()
		case <-ctx.Done():
			return ctx.Err()
		}

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), s.cfg.Server.ShutdownTimeout)
		defer shutdownCancel()

		if err := s.app.ShutdownWithContext(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}

		s.logger.Info("server stopped gracefully")
		return nil
	})

	return g.Wait()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
