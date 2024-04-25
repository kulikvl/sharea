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

		availSpace, err := s.Storage.CalculateAvailableSpace()
		if err != nil {
			availSpace = 0
		}

		return c.Render("index", fiber.Map{
			"Name":              hostname,
			"IP":                s.IP,
			"CapacityAvailable": availSpace,
			"Capacity":          s.Storage.Capacity,
		})
	})

	s.App.Get("/test", func(c *fiber.Ctx) error {
		return c.Render("test", fiber.Map{})
	})
}
