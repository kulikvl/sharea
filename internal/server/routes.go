package server

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

func (s *Server) setupRoutes() {
	s.App.Static("/static", "./web/static")

	s.App.Get("/", func(c *fiber.Ctx) error {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		return c.Render("index", fiber.Map{
			"Name":     hostname,
			"IP":       s.IP,
			"Capacity": s.Storage.Capacity,
		})
	})
}
