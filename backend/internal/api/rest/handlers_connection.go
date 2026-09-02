package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (hc *HandlerContext) listConnections(c *fiber.Ctx) error {
	userID := getUserID(c)
	zap.L().Info("list connections", zap.String("user_id", userID))

	return c.JSON([]fiber.Map{})
}

func (hc *HandlerContext) getActiveConnection(c *fiber.Ctx) error {
	userID := getUserID(c)
	zap.L().Info("get active connection", zap.String("user_id", userID))

	return c.JSON(fiber.Map{
		"status":     "disconnected",
		"protocol":   nil,
		"server":     nil,
		"connected_at": nil,
		"bytes_rx":   0,
		"bytes_tx":   0,
		"latency_ms": nil,
	})
}

type startConnectionRequest struct {
	ServerID string `json:"server_id"`
	Protocol string `json:"protocol"`
}

func (hc *HandlerContext) startConnection(c *fiber.Ctx) error {
	userID := getUserID(c)

	var req startConnectionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	connectionID := uuid.New()

	zap.L().Info("connection started",
		zap.String("user_id", userID),
		zap.String("connection_id", connectionID.String()),
		zap.String("server_id", req.ServerID),
		zap.String("protocol", req.Protocol),
	)

	return c.JSON(fiber.Map{
		"connection_id": connectionID.String(),
		"status":        "connected",
		"started_at":    time.Now(),
	})
}

func (hc *HandlerContext) endConnection(c *fiber.Ctx) error {
	userID := getUserID(c)
	zap.L().Info("connection ended", zap.String("user_id", userID))

	return c.JSON(fiber.Map{
		"status": "disconnected",
	})
}

func (hc *HandlerContext) getConnectionStats(c *fiber.Ctx) error {
	userID := getUserID(c)
	zap.L().Info("connection stats", zap.String("user_id", userID))

	return c.JSON(fiber.Map{
		"total_bytes_rx": 0,
		"total_bytes_tx": 0,
		"total_sessions": 0,
		"total_duration": "0s",
	})
}
