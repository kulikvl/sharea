package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/kulikvl/sharea/internal/storage"
	"github.com/kulikvl/sharea/internal/utils"
	"log"
	"net/url"
)

type Server struct {
	Port    int
	IP      string
	Storage storage.Storage
	App     *fiber.App
}

func New(port int, path string, capacity int64) (*Server, error) {
	ip, err := utils.GetLocalIp()
	if err != nil {
		return nil, fmt.Errorf("failed to get server's local ip")
	}

	return &Server{
		Port: port,
		IP:   ip,
		Storage: storage.Storage{
			Path:     path,
			Capacity: capacity,
		},
		App: nil,
	}, nil
}

func (s *Server) URL() string {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", s.IP, s.Port),
	}
	return u.String()
}

func (s *Server) Run() {
	// init the template engine (Go's built-in engine to render HTML pages)
	engine := html.New("./web", ".html")

	s.App = fiber.New(fiber.Config{
		Views:                 engine,
		BodyLimit:             1 * 1024 * 1024 * 1024, // 1 GB
		DisableStartupMessage: false,
	})

	s.App.Use(logger.New())

	s.setupRoutes()
	s.setupApi()
	s.setupWebSockets()

	log.Fatal(s.App.Listen(fmt.Sprintf(":%d", s.Port)))
}
