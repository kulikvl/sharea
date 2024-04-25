package server

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	fiberUtils "github.com/gofiber/fiber/v2/utils"
	"io"
	"os"
)

func (s *Server) setupApi() {
	api := s.App.Group("/api")

	// Enable request body streaming to read uploaded files in chunks, rather than store the whole file in RAM
	// https://github.com/gofiber/recipes/blob/master/stream-request-body/main.go
	s.App.Server().StreamRequestBody = true

	api.Post("/upload/:filename", func(c *fiber.Ctx) error {
		filename := fiberUtils.CopyString(c.Params("filename"))
		filesize := c.Context().Request.Header.ContentLength()

		// check if file with the same name already exists in the storage
		// todo: delete that check, instead add automatic renaming
		if fileExists, err := s.Storage.FileExists(filename); fileExists || err != nil {
			if err != nil {
				return fmt.Errorf("failed to get info about storage files: %w", err)
			}
			return fmt.Errorf("file with name %s already exists in the storage", filename)
		}

		// check for available space
		availableSpace, err := s.Storage.CalculateAvailableSpace()
		if err != nil {
			return fmt.Errorf("failed to calculate available storage space: %w", err)
		}
		if int64(filesize) > availableSpace {
			return c.Status(fiber.StatusRequestEntityTooLarge).SendString(fmt.Sprintf("total file size exceeds the storage capacity (%d bytes)", s.Storage.Capacity))
		}

		file, err := os.Create(fmt.Sprintf("%s/%s", s.Storage.Path, filename))
		if err != nil {
			return fmt.Errorf("failed to create upload file %s: %w", filename, err)
		}

		reader := c.Context().RequestBodyStream()
		// todo: make that separate configurable variable such as "uploadBufferSize"
		buffer := make([]byte, 0, 50*1024*1024) // 50 MiB / s

		for {
			length, err := io.ReadFull(reader, buffer[:cap(buffer)])
			buffer = buffer[:length]

			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				if !errors.Is(err, io.ErrUnexpectedEOF) {
					return fmt.Errorf("failed to read body stream: %w", err)
				}
			}

			if _, err := file.Write(buffer); err != nil {
				return fmt.Errorf("failed to write %d bytes to file %s: %w", length, filename, err)
			}
		}

		if err := file.Sync(); err != nil {
			return fmt.Errorf("failed to flush (sync) file %s: %w", filename, err)
		}

		return c.Status(fiber.StatusAccepted).SendString(fmt.Sprintf("File %s was successfully uploaded", filename))
	})

	api.Get("/download/:filename", func(c *fiber.Ctx) error {
		filename := fiberUtils.CopyString(c.Params("filename"))
		return c.Download(fmt.Sprintf("%s/%s", s.Storage.Path, filename))
	})
}
