package docs

import (
	_ "embed"
	"strings"

	"github.com/gofiber/fiber/v3"
)

//go:embed index.html
var indexHTML []byte

//go:embed openapi.yml
var openAPISpec []byte

func RegisterRoutes(app *fiber.App) {
	app.Get("/docs/openapi.yml", func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "application/yaml")
		return c.Send(openAPISpec)
	})

	app.Get("/docs", func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.Send(indexHTML)
	})

	app.Get("/docs/*", func(c fiber.Ctx) error {
		path := strings.TrimPrefix(c.Path(), "/docs/")
		if path == "" || path == "/" {
			c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
			return c.Send(indexHTML)
		}
		return c.Status(fiber.StatusNotFound).SendString("not found")
	})
}
