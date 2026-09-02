package rest

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/nebula-vpn/backend/internal/auth"
)

type contextKey string

const (
	ContextUserID  contextKey = "user_id"
	ContextEmail   contextKey = "email"
	ContextTokenID contextKey = "token_id"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	sqlInjection  = regexp.MustCompile(`(?i)(\b(SELECT|INSERT|UPDATE|DELETE|DROP|UNION|ALTER|CREATE|TRUNCATE|EXEC|EXECUTE)\b)`)
)

func authMiddleware(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		claims, err := jwtManager.ValidateAccessToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		c.Locals(string(ContextUserID), claims.UserID)
		c.Locals(string(ContextEmail), claims.Email)
		c.Locals(string(ContextTokenID), claims.ID)

		return c.Next()
	}
}

func getUserID(c *fiber.Ctx) string {
	if v := c.Locals(string(ContextUserID)); v != nil {
		return v.(string)
	}
	return ""
}

func getUserEmail(c *fiber.Ctx) string {
	if v := c.Locals(string(ContextEmail)); v != nil {
		return v.(string)
	}
	return ""
}

// Input validation
func validateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func validatePassword(pw string) bool {
	if len(pw) < 8 || len(pw) > 128 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pw)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pw)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(pw)
	return hasUpper && hasLower && hasDigit
}

func sanitizeInput(s string) string {
	return sqlInjection.ReplaceAllString(s, "")
}

func validateRequestBody(c *fiber.Ctx, requiredFields ...string) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	for _, field := range requiredFields {
		if _, ok := body[field]; !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing required field: " + field,
			})
		}
	}
	return c.Next()
}
