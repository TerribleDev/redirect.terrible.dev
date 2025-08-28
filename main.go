package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

// Hardcoded redirect rules
var (
	// Host-based redirects
	hostRedirects = map[string]string{
		"mail.terrible.dev": "https://mail.tommyparnell.com",
	}

	// Path-based redirects
	pathRedirects = map[string]string{
		"/test": "https://blog.terrible.dev",
	}
)

func redirectHandler(c *fiber.Ctx) error {
	// Get the actual host, checking X-Forwarded-Host first
	host := c.Get("X-Forwarded-Host")
	if host == "" {
		host = c.Hostname()
	}

	// Check for host-based redirects first
	if redirectURL, exists := hostRedirects[host]; exists {
		return c.Redirect(redirectURL, fiber.StatusMovedPermanently)
	}

	// Get the actual path, checking X-Forwarded-Uri first
	path := c.Get("X-Forwarded-Uri")
	if path == "" {
		path = c.Path()
	}

	// Check for path-based redirects
	if redirectURL, exists := pathRedirects[path]; exists {
		return c.Redirect(redirectURL, fiber.StatusMovedPermanently)
	}

	// If no redirect rule matches, return 404
	return c.SendStatus(fiber.StatusNotFound)
}

func main() {
	// Create Fiber app with minimal configuration for lowest memory usage
	app := fiber.New(fiber.Config{
		Prefork:          false,
		DisableKeepalive: false,
		ServerHeader:     "",
		AppName:          "redirect",
	})

	// Use a single handler function to minimize memory overhead
	app.All("*", redirectHandler)

	// Get port from environment variable (Heroku sets this)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port for local development
	}

	// Start server on the specified port
	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}
