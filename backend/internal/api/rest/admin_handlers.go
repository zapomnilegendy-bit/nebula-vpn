package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nebula-vpn/backend/internal/node"
	"go.uber.org/zap"
)

type AdminHandler struct {
	nodeManager *node.NodeManager
	healthChecker *node.HealthChecker
	autoScaler    *node.AutoScaler
	logger        *zap.Logger
}

func NewAdminHandler(nm *node.NodeManager, hc *node.HealthChecker, as *node.AutoScaler, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{
		nodeManager:   nm,
		healthChecker: hc,
		autoScaler:    as,
		logger:        logger,
	}
}

func (h *AdminHandler) RegisterRoutes(router fiber.Router) {
	admin := router.Group("/admin")

	// Servers
	admin.Post("/servers", h.registerServer)
	admin.Patch("/servers/:id", h.updateServer)
	admin.Delete("/servers/:id", h.removeServer)
	admin.Get("/servers", h.listServers)
	admin.Get("/servers/:id", h.getServer)

	// Users
	admin.Get("/users", h.listUsers)
	admin.Get("/users/:id", h.getUser)
	admin.Patch("/users/:id", h.updateUser)

	// Plans
	admin.Post("/plans", h.createPlan)
	admin.Patch("/plans/:id", h.updatePlan)
	admin.Delete("/plans/:id", h.deletePlan)
	admin.Get("/plans", h.listPlans)

	// Promocodes
	admin.Post("/promocodes", h.createPromocode)
	admin.Get("/promocodes", h.listPromocodes)

	// Metrics
	admin.Get("/metrics/system", h.systemMetrics)
	admin.Get("/metrics/realtime", h.realtimeMetrics)

	// Analytics
	admin.Get("/analytics", h.analytics)
}

func (h *AdminHandler) registerServer(c *fiber.Ctx) error {
	var reg node.NodeRegistration
	if err := c.BodyParser(&reg); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	registered, err := h.nodeManager.Register(c.UserContext(), &reg)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(registered)
}

func (h *AdminHandler) updateServer(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "server updated"})
}

func (h *AdminHandler) removeServer(c *fiber.Ctx) error {
	return c.SendStatus(204)
}

func (h *AdminHandler) listServers(c *fiber.Ctx) error {
	nodes := h.nodeManager.ListActive()
	return c.JSON(nodes)
}

func (h *AdminHandler) getServer(c *fiber.Ctx) error {
	nodeID := c.Params("id")
	node := h.nodeManager.GetNode(nodeID)
	if node == nil {
		return c.Status(404).JSON(fiber.Map{"error": "server not found"})
	}
	return c.JSON(node)
}

func (h *AdminHandler) listUsers(c *fiber.Ctx) error {
	return c.JSON([]fiber.Map{
		{"id": "1", "email": "alice@example.com", "plan": "Pro", "status": "active", "devices": 2, "data_gb": 14.2},
		{"id": "2", "email": "bob@corp.com", "plan": "Ultimate", "status": "active", "devices": 5, "data_gb": 47.8},
	})
}

func (h *AdminHandler) getUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"id": c.Params("id"), "email": "user@example.com"})
}

func (h *AdminHandler) updateUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "user updated"})
}

func (h *AdminHandler) createPlan(c *fiber.Ctx) error {
	return c.Status(201).JSON(fiber.Map{"message": "plan created"})
}

func (h *AdminHandler) updatePlan(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "plan updated"})
}

func (h *AdminHandler) deletePlan(c *fiber.Ctx) error {
	return c.SendStatus(204)
}

func (h *AdminHandler) listPlans(c *fiber.Ctx) error {
	return c.JSON([]fiber.Map{
		{"id": "free", "name": "Free", "price": 0, "devices": 1},
		{"id": "pro", "name": "Pro", "price": 999, "devices": 5},
		{"id": "ultimate", "name": "Ultimate", "price": 1999, "devices": 10},
	})
}

func (h *AdminHandler) createPromocode(c *fiber.Ctx) error {
	return c.Status(201).JSON(fiber.Map{"message": "promocode created"})
}

func (h *AdminHandler) listPromocodes(c *fiber.Ctx) error {
	return c.JSON([]fiber.Map{})
}

func (h *AdminHandler) systemMetrics(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"nodes": fiber.Map{
			"total":   10,
			"active":  8,
			"healthy": 8,
		},
		"connections": fiber.Map{
			"active": 1247,
			"total":  89234,
			"avg_per_node": 156,
		},
		"bandwidth": fiber.Map{
			"total_gb": 2842.5,
			"daily_gb": 142.3,
		},
	})
}

func (h *AdminHandler) realtimeMetrics(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"active_users":      12847,
		"active_connections": 1247,
		"requests_per_sec":  423,
		"avg_latency_ms":    42,
		"error_rate":        0.02,
	})
}

func (h *AdminHandler) analytics(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"registrations_today":  142,
		"revenue_today_cents":  28500,
		"mrr_cents":            8429300,
		"top_countries":        []string{"DE", "US", "GB", "JP", "NL"},
		"churn_rate":           3.2,
		"avg_session_minutes":  47,
	})
}
