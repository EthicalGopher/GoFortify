package server

import (
	"bytes"
	"fmt"
	"github.com/EthicalGopher/SentinelShield/vulnerabilities"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

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
func Server() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()

		ip := c.IP()
		method := c.Method()
		path := c.Path()

		// -------- Query Params --------
		query := c.Queries()
		if len(query) > 0 {
			if vulnerabilities.SqlInjection(c, query) {
				return c.Status(fiber.StatusForbidden).
					SendString("Blocked by SentinelShield (SQL Injection)")
			}
		}

		// -------- Allow & Forward --------

		latency := time.Since(start)

		// -------- Logging --------
		log.Printf(
			"[%s] %s | %s %s | Latency: %s",
			time.Now().Format("2006/01/02 15:04:05"),
			ip,
			method,
			path,
			latency,
		)

		if len(query) > 0 {
			fmt.Print("\tQueries:")
			for k, v := range query {
				fmt.Printf(" %s=%s", k, v)
			}
			fmt.Println()
		}

		// -------- Body (multipart only) --------
		ct := c.Get("Content-Type")
		if strings.HasPrefix(ct, "multipart/form-data") {
			form, err := c.MultipartForm()
			if err == nil && form != nil && len(form.Value) > 0 {
				if vulnerabilities.SqlInjectionBody(c, form.Value) {
					return c.Status(fiber.StatusForbidden).
						SendString("Blocked by SentinelShield (SQL Injection)")
				}
			}
			if form.Value != nil {
				fmt.Print("\t\tBody : ")
				for key, values := range form.Value {
					for _, v := range values {
						fmt.Print("\t"+key, "->", v)
					}
				}
				fmt.Println()
			}

		}

		return ForwardToBackend(c, "http://localhost:8080")
	})

	log.Println("🛡 SentinelShield (Fiber) Proxy STARTED on :5174")
	log.Fatal(app.Listen(":5174"))
}
