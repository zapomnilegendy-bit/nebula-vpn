package rest

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (hc *HandlerContext) getSubscription(c *fiber.Ctx) error {
	userID := getUserID(c)
	zap.L().Info("get subscription", zap.String("user_id", userID))

	return c.JSON(fiber.Map{
		"plan":      "free",
		"status":    "active",
		"expires_at": nil,
		"max_devices": 1,
	})
}

type createSubscriptionRequest struct {
	PlanID      string `json:"plan_id"`
	PaymentMethod string `json:"payment_method"`
}

func (hc *HandlerContext) createSubscription(c *fiber.Ctx) error {
	userID := getUserID(c)

	var req createSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	zap.L().Info("create subscription",
		zap.String("user_id", userID),
		zap.String("plan_id", req.PlanID),
	)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"plan":      req.PlanID,
		"status":    "active",
		"message":   "subscription created",
	})
}

func (hc *HandlerContext) updateSubscription(c *fiber.Ctx) error {
	subID := c.Params("id")

	return c.JSON(fiber.Map{
		"id":      subID,
		"status":  "updated",
		"message": "subscription updated",
	})
}

func (hc *HandlerContext) cancelSubscription(c *fiber.Ctx) error {
	subID := c.Params("id")
	zap.L().Info("cancel subscription", zap.String("subscription_id", subID))

	return c.JSON(fiber.Map{
		"id":      subID,
		"status":  "canceled",
		"message": "subscription canceled",
	})
}

func (hc *HandlerContext) listPlans(c *fiber.Ctx) error {
	return c.JSON([]fiber.Map{
		{
			"id":               "plan_free",
			"name":             "Free",
			"price_cents":      0,
			"currency":         "USD",
			"max_devices":      1,
			"max_speed_mbps":   50,
			"protocols_enabled": []string{"wireguard"},
		},
		{
			"id":               "plan_pro",
			"name":             "Pro",
			"price_cents":      999,
			"currency":         "USD",
			"max_devices":      5,
			"max_speed_mbps":   500,
			"protocols_enabled": []string{"wireguard", "vless_reality", "hysteria2"},
		},
		{
			"id":               "plan_ultimate",
			"name":             "Ultimate",
			"price_cents":      1999,
			"currency":         "USD",
			"max_devices":      10,
			"max_speed_mbps":   0,
			"protocols_enabled": []string{"wireguard", "vless_reality", "hysteria2", "tuic"},
		},
	})
}

type redeemPromocodeRequest struct {
	Code string `json:"code"`
}

func (hc *HandlerContext) redeemPromocode(c *fiber.Ctx) error {
	var req redeemPromocodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	return c.JSON(fiber.Map{
		"code":    req.Code,
		"discount": 20,
		"message": "promocode applied",
	})
}
