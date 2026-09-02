package rest

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type serverResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Country   string   `json:"country"`
	City      string   `json:"city"`
	Protocols []string `json:"protocols"`
	Load      int      `json:"load_percent"`
	Latency   int      `json:"latency_ms"`
	IsHealthy bool     `json:"is_healthy"`
	IsActive  bool     `json:"is_active"`
}

func (hc *HandlerContext) listServers(c *fiber.Ctx) error {
	country := c.Query("country")

	rows, err := hc.DB.Query(c.UserContext(), `
		SELECT id::text, name, country, COALESCE(city, ''), load_percent, is_healthy, is_active
		FROM vpn_servers
		WHERE ($1 = '' OR $1 IS NULL OR country = $1)
		ORDER BY country, city
	`, country)
	if err != nil {
		zap.L().Error("failed to query servers", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	defer rows.Close()

	servers := make([]serverResponse, 0)
	for rows.Next() {
		var s serverResponse

		if err := rows.Scan(&s.ID, &s.Name, &s.Country, &s.City, &s.Load, &s.IsHealthy, &s.IsActive); err != nil {
			zap.L().Error("failed to scan server row", zap.Error(err))
			continue
		}

		s.Protocols = []string{"wireguard", "vless_reality"}
		if s.IsHealthy && s.IsActive {
			s.Latency = 15 + (s.Load / 5)
		}

		servers = append(servers, s)
	}

	return c.JSON(servers)
}

func (hc *HandlerContext) getServer(c *fiber.Ctx) error {
	serverID := c.Params("id")

	var s serverResponse

	err := hc.DB.QueryRow(c.UserContext(), `
		SELECT id::text, name, country, COALESCE(city, ''), load_percent, is_healthy, is_active
		FROM vpn_servers WHERE id::text = $1
	`, serverID).Scan(&s.ID, &s.Name, &s.Country, &s.City, &s.Load, &s.IsHealthy, &s.IsActive)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "server not found"})
	}

	s.Protocols = []string{"wireguard", "vless_reality"}
	if s.IsHealthy && s.IsActive {
		s.Latency = 15 + (s.Load / 5)
	}

	return c.JSON(s)
}

func (hc *HandlerContext) recommendServer(c *fiber.Ctx) error {
	var s serverResponse

	err := hc.DB.QueryRow(c.UserContext(), `
		SELECT id::text, name, country, COALESCE(city, ''), load_percent, is_healthy, is_active
		FROM vpn_servers
		WHERE is_healthy = true AND is_active = true
		ORDER BY load_percent ASC
		LIMIT 1
	`).Scan(&s.ID, &s.Name, &s.Country, &s.City, &s.Load, &s.IsHealthy, &s.IsActive)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "no healthy servers"})
	}

	s.Protocols = []string{"wireguard", "vless_reality"}
	s.Latency = 15 + (s.Load / 5)

	return c.JSON(s)
}
