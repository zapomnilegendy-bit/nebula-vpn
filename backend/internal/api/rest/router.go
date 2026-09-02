package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nebula-vpn/backend/internal/auth"
	"github.com/nebula-vpn/backend/internal/config"
)

type HandlerContext struct {
	JWTManager   *auth.JWTManager
	RateLimiter  *auth.RateLimiter
	Config       *config.Config
	DB           *pgxpool.Pool
	AdminHandler *AdminHandler
}

func NewRouter(hc *HandlerContext) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "NebulaVPN API",
		ReadTimeout:           hc.Config.Server.ReadTimeout,
		WriteTimeout:          hc.Config.Server.WriteTimeout,
		IdleTimeout:           120 * time.Second,
		Prefork:               false,
		ServerHeader:          "NebulaVPN",
		DisableStartupMessage: false,
		Concurrency:           262144,
		ReadBufferSize:        4096,
		WriteBufferSize:       4096,
		ReduceMemoryUsage:     true,
		CaseSensitive:         true,
		StrictRouting:         true,
		EnableTrustedProxyCheck: true,
		TrustedProxies:        []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
	})

	// Global middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip}:${port} ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "2006-01-02T15:04:05Z07:00",
		TimeZone:   "UTC",
	}))

	// Performance: compression (level 2 = balanced)
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Security: CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://app.nebula-vpn.io, https://admin.nebula-vpn.io",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-Id",
		ExposeHeaders:    "X-Request-Id, X-RateLimit-Remaining",
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Security: rate limiter (distributed when Redis available)
	if hc.Config.RateLimit.Enabled {
		app.Use(limiter.New(limiter.Config{
			Max:        hc.Config.RateLimit.RequestsPerSec,
			Expiration: 1 * time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "rate limit exceeded",
					"retry_after": 60,
				})
			},
		}))
	}

	// Security: security headers
	app.Use(securityHeaders)

	// Metrics: Prometheus /metrics
	app.Get("/metrics", monitor.New(monitor.Config{
		Title:   "NebulaVPN Metrics",
		Refresh: 10,
	}))

	// Pprof: performance profiling (disabled in production)
	if hc.Config.Server.EnablePprof {
		setupPprof(app)
	}

	// Health
	app.Get("/health", healthCheck)
	app.Get("/ready", readinessCheck)

	// API v1
	v1 := app.Group("/v1")

	// Auth (no auth required, but with stricter rate limit)
	authGroup := v1.Group("/auth")
	authGroup.Post("/register", hc.registerUser)
	authGroup.Post("/login", loginRateLimiter(), hc.loginUser)
	authGroup.Post("/refresh", hc.refreshToken)
	authGroup.Post("/logout", hc.logoutUser)

	// Protected routes
	protected := v1.Group("")
	protected.Use(authMiddleware(hc.JWTManager))

	// Devices
	devices := protected.Group("/devices")
	devices.Get("/", hc.listDevices)
	devices.Post("/", hc.registerDevice)
	devices.Get("/:id", hc.getDevice)
	devices.Patch("/:id", hc.updateDevice)
	devices.Delete("/:id", hc.removeDevice)
	devices.Get("/:id/config", hc.getDeviceConfig)

	// Servers
	servers := protected.Group("/servers")
	servers.Get("/", hc.listServers)
	servers.Get("/:id", hc.getServer)
	servers.Get("/recommend", hc.recommendServer)

	// Connections
	connections := protected.Group("/connections")
	connections.Get("/", hc.listConnections)
	connections.Get("/active", hc.getActiveConnection)
	connections.Post("/start", hc.startConnection)
	connections.Post("/end", hc.endConnection)
	connections.Get("/stats", hc.getConnectionStats)

	// Subscriptions
	subscriptions := protected.Group("/subscriptions")
	subscriptions.Get("/", hc.getSubscription)
	subscriptions.Post("/", hc.createSubscription)
	subscriptions.Patch("/:id", hc.updateSubscription)
	subscriptions.Delete("/:id", hc.cancelSubscription)

	// Plans
	protected.Get("/plans", hc.listPlans)

	// Promocodes
	protected.Post("/promocodes/redeem", hc.redeemPromocode)

	// Admin
	if hc.AdminHandler != nil {
		hc.AdminHandler.RegisterRoutes(protected)
	}

	return app
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "nebula-vpn-api",
	})
}

func readinessCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ready",
		"service": "nebula-vpn-api",
	})
}

func securityHeaders(c *fiber.Ctx) error {
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-XSS-Protection", "1; mode=block")
	c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(self)")
	c.Set("Cross-Origin-Resource-Policy", "same-origin")
	c.Set("Cross-Origin-Opener-Policy", "same-origin")
	return c.Next()
}

func loginRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "login:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":      "too many login attempts",
				"retry_after": 60,
			})
		},
	})
}
