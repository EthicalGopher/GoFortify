// Package server implements the core reverse proxy and security filtering logic.
// It intercepts incoming requests, applies rate limiting, evaluates vulnerabilities,
// and forwards legitimate traffic to the backend server.
package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/EthicalGopher/GoFortify/shared"
	"github.com/EthicalGopher/GoFortify/vulnerabilities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

var (
	logger = shared.NewLogger()
)

const (
	reset  = "\033[0m"
	gray   = "\033[90m"
	cyan   = "\033[36m"
	yellow = "\033[33m"
	green  = "\033[32m"
	red    = "\033[31m"
)

// ForwardToBackend proxies an incoming HTTP request to the designated backend server.
// It replicates the method, body, and headers from the original Fiber context, executes
// the HTTP request against the backend, and writes the response back to the client.
func ForwardToBackend(c *fiber.Ctx, backendBaseURL string) error {
	url := backendBaseURL + c.OriginalURL()

	req, err := http.NewRequest(
		c.Method(),
		url,
		bytes.NewReader(c.Body()),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString("Failed to create backend request")
	}

	c.Request().Header.VisitAll(func(k, v []byte) {
		req.Header.Set(string(k), string(v))
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).
			SendString("Backend server not reachable")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			SendString("Failed to read backend response")
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	return c.Status(resp.StatusCode).Send(respBody)
}
// Server initializes and starts the Fiber application that acts as the security proxy.
// It sets up CORS, applies rate limiting, and configures the main middleware to inspect
// query parameters and request bodies for SQL Injection and XSS before forwarding.
func Server(backendURL string, proxyPort int, filename_ratelimit string, filename_xss string, filename_sql string, ratelimit int) error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        ratelimit,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			vulnerabilities.LogRateLimit(c, filename_ratelimit)
			return c.Status(429).SendString("Too many requests - blocked by GoFortify")
		},
	}))
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()

		ip := c.IP()
		method := c.Method()
		path := c.Path()

		latency := time.Since(start)

		// -------- Logging --------
		logger.Println(gray + "--------------------------------------------------" + reset)
		logger.Printf(
			cyan+"Time     :"+reset+" %s\n"+
				cyan+"IP       :"+reset+" %s\n"+
				cyan+"Request  :"+reset+" %s %s\n"+
				yellow+"Latency  :"+reset+" %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			ip,
			method,
			path,
			latency,
		)

		// -------- Query Params --------
		query := c.Queries()
		if len(query) > 0 {
			if vulnerabilities.SqlInjection(c, query, filename_sql) {
				logger.Println(red + "Blocked by GoFortify (SQL Injection in Query)" + reset)
				return c.Status(fiber.StatusForbidden).
					SendString("Blocked by GoFortify (SQL Injection)")
			}
			if vulnerabilities.XSSInjection(c, query, filename_xss) {
				logger.Println(red + "Blocked by GoFortify (XSS attempt in Query)" + reset)
				return c.Status(fiber.StatusForbidden).SendString("Blocked by GoFortify due to XSS attempt")
			}

			logger.Println(green + "Query Parameters:" + reset)
			for k, v := range query {
				logger.Printf("  %s- %s = %s%s\n", green, k, v, reset)
			}
		}

		// -------- Body (multiple content types) --------
		ct := c.Get("Content-Type")
		if strings.HasPrefix(ct, "multipart/form-data") || strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			// Fiber's BodyParser handles both multipart and urlencoded for map[string][]string if we use it correctly,
			// but c.MultipartForm() is specific to multipart.
			// Let's use c.BodyParser to be more generic if possible, or handle them separately.

			if strings.HasPrefix(ct, "multipart/form-data") {
				form, err := c.MultipartForm()
				if err == nil && form != nil && len(form.Value) > 0 {
					if vulnerabilities.SqlInjectionBody(c, form.Value, filename_sql) {
						logger.Println(red + "Blocked by GoFortify (SQL Injection in Multipart Body)" + reset)
						return c.Status(fiber.StatusForbidden).
							SendString("Blocked by GoFortify (SQL Injection)")
					}
					if vulnerabilities.XSSInjectionBody(c, form.Value, filename_xss) {
						logger.Println(red + "Blocked by GoFortify (XSS attempt in Multipart Body)" + reset)
						return c.Status(fiber.StatusForbidden).SendString("Blocked by GoFortify due to XSS attempt")
					}
					logger.Println(yellow + "Multipart Body:" + reset)
					for key, values := range form.Value {
						for _, v := range values {
							logger.Printf("  %s- %s -> %s%s\n", green, key, v, reset)
						}
					}
				}
			} else {
				// application/x-www-form-urlencoded
				bodyStr := string(c.Body())
				if values, err := url.ParseQuery(bodyStr); err == nil && len(values) > 0 {
					bodyMap := make(map[string]string)
					for k, v := range values {
						if len(v) > 0 {
							bodyMap[k] = v[0]
						}
					}
					if vulnerabilities.SqlInjection(c, bodyMap, filename_sql) {
						logger.Println(red + "Blocked by GoFortify (SQL Injection in URL-Encoded Body)" + reset)
						return c.Status(fiber.StatusForbidden).
							SendString("Blocked by GoFortify (SQL Injection)")
					}
					if vulnerabilities.XSSInjection(c, bodyMap, filename_xss) {
						logger.Println(red + "Blocked by GoFortify (XSS attempt in URL-Encoded Body)" + reset)
						return c.Status(fiber.StatusForbidden).SendString("Blocked by GoFortify due to XSS attempt")
					}
					logger.Println(yellow + "URL-Encoded Body:" + reset)
					for k, v := range bodyMap {
						logger.Printf("  %s- %s = %s%s\n", green, k, v, reset)
					}
				}
			}
		} else if strings.HasPrefix(ct, "application/json") {
			var body map[string]interface{}
			if err := c.BodyParser(&body); err == nil && len(body) > 0 {
				if vulnerabilities.SqlInjectionInterface(c, body, filename_sql) {
					logger.Println(red + "Blocked by GoFortify (SQL Injection in JSON Body)" + reset)
					return c.Status(fiber.StatusForbidden).
						SendString("Blocked by GoFortify (SQL Injection)")
				}
				if vulnerabilities.XSSInjectionInterface(c, body, filename_xss) {
					logger.Println(red + "Blocked by GoFortify (XSS attempt in JSON Body)" + reset)
					return c.Status(fiber.StatusForbidden).SendString("Blocked by GoFortify due to XSS attempt")
				}
				logger.Println(yellow + "JSON Body:" + reset)
				for k, v := range body {
					logger.Printf("  %s- %s = %v%s\n", green, k, v, reset)
				}
			}
		}

		return ForwardToBackend(c, backendURL)
	})

	logger.Printf("GoFortify Proxy STARTED on :%d\n", proxyPort)
	return app.Listen(fmt.Sprintf(":%d", proxyPort))
}
