package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Hardcoded redirect rules
var (
	// Host-based redirects
	hostRedirects = map[string]string{
		"mail.terrible.dev":     "https://mail.tommyparnell.com/mail",
		"calendar.terrible.dev": "https://mail.tommyparnell.com/cloud/calendar",
		"cal.terrible.dev":      "https://cal.com/terribledev",
	}

	// Path-based redirects
	pathRedirects = map[string]string{
		"/test": "https://blog.terrible.dev",
	}
)

var listHosts = map[string]bool{
	"aka.terrible.dev":      true,
	"redirect.terrible.dev": true,
}

func listRoutesHandler(c *fiber.Ctx) error {
	host := c.Get("X-Forwarded-Host")
	if host == "" {
		host = c.Hostname()
	}
	if !listHosts[host] {
		return redirectHandler(c)
	}

	var sb strings.Builder
	sb.WriteString("Redirect Routes\n\n")

	sb.WriteString("Host-based redirects:\n")
	hosts := make([]string, 0, len(hostRedirects))
	for h := range hostRedirects {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)
	for _, h := range hosts {
		sb.WriteString(fmt.Sprintf("  %s -> %s\n", h, hostRedirects[h]))
	}

	sb.WriteString("\nPath-based redirects:\n")
	paths := make([]string, 0, len(pathRedirects))
	for p := range pathRedirects {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		sb.WriteString(fmt.Sprintf("  %s -> %s\n", p, pathRedirects[p]))
	}

	c.Set("Content-Type", "text/plain")
	return c.SendString(sb.String())
}

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

	// List all routes on /
	app.Get("/", listRoutesHandler)

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
