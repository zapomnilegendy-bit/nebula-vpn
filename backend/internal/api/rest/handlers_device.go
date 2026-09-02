package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type deviceRequest struct {
	Name     string `json:"name" validate:"required"`
	Platform string `json:"platform" validate:"required,oneof=macos windows ios android"`
}

type deviceResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Platform  string `json:"platform"`
	PublicKey string `json:"public_key"`
	Address   string `json:"address"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

func (hc *HandlerContext) listDevices(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// TODO: fetch from DB
	zap.L().Info("list devices", zap.String("user_id", userID))

	return c.JSON([]deviceResponse{})
}

func (hc *HandlerContext) registerDevice(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req deviceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Name == "" || req.Platform == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name and platform are required",
		})
	}

	deviceID := uuid.New()

	zap.L().Info("device registered",
		zap.String("user_id", userID),
		zap.String("device_id", deviceID.String()),
		zap.String("platform", req.Platform),
	)

	return c.Status(fiber.StatusCreated).JSON(deviceResponse{
		ID:        deviceID.String(),
		Name:      req.Name,
		Platform:  req.Platform,
		IsActive:  true,
		CreatedAt: "now",
	})
}

func (hc *HandlerContext) getDevice(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	deviceID := c.Params("id")
	zap.L().Info("get device", zap.String("user_id", userID), zap.String("device_id", deviceID))

	// TODO: fetch from DB
	return c.JSON(deviceResponse{
		ID:        deviceID,
		Name:      "My Device",
		Platform:  "macos",
		IsActive:  true,
		CreatedAt: "now",
	})
}

func (hc *HandlerContext) updateDevice(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	return c.JSON(fiber.Map{
		"message": "device updated",
	})
}

func (hc *HandlerContext) removeDevice(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (hc *HandlerContext) getDeviceConfig(c *fiber.Ctx) error {
	deviceID := c.Params("id")

	// Return WireGuard config
	config := `[Interface]
PrivateKey = <generated>
Address = 10.0.0.2/32
DNS = 1.1.1.1, 1.0.0.1
MTU = 1420

[Peer]
PublicKey = <server-public-key>
Endpoint = vpn.nebula-vpn.io:51820
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 25`

	zap.L().Info("config requested", zap.String("device_id", deviceID))

	return c.JSON(fiber.Map{
		"config": config,
	})
}
