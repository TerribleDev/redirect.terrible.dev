# URL Redirect Service

A minimal Go application for handling URL redirects with the smallest possible memory footprint.

## Features

- **Host-based redirects**: Redirects based on the incoming host header
- **Path-based redirects**: Redirects based on the URL path
- **Proxy support**: Handles `X-Forwarded-Host` and `X-Forwarded-Uri` headers from reverse proxies
- **Protocol agnostic**: Works with both HTTP and HTTPS requests
- **Minimal memory usage**: Uses only Go's standard library with no external dependencies

## Current Redirect Rules

### Host-based
- `mail.terrible.dev` → `https://mail.tommyparnell.com`

### Path-based  
- `/test` → `https://blog.terrible.dev`

## Usage

1. Run the application:
   ```bash
   go run main.go
   ```

2. The server will start on port 8080

3. Test the redirects:
   ```bash
   # Test host-based redirect
   curl -H "Host: mail.terrible.dev" http://localhost:8080
   
   # Test path-based redirect
   curl http://localhost:8080/test
   
   # Test with X-Forwarded headers (simulating a reverse proxy)
   curl -H "X-Forwarded-Host: mail.terrible.dev" http://localhost:8080
   curl -H "X-Forwarded-Uri: /test" http://localhost:8080
   ```

## Building

To build a standalone binary:
```bash
go build -o redirect main.go
```

## Memory Optimization

This application is designed for minimal memory usage:
- Uses only Go's standard `net/http` package
- Single handler function for all routes
- Simple map lookups for redirect rules
- No external dependencies or frameworks
- Supports reverse proxy headers without additional overhead
