package vulnerabilities

import (
	"encoding/json"
	"os"
	"time"

	"github.com/EthicalGopher/GoFortify/shared"
	"github.com/gofiber/fiber/v2"
)

func LogRateLimit(c *fiber.Ctx, rateLimitFile string) {
	logEntry := shared.RateLimitLog{
		Ip:     c.IP(),
		Path:   c.Path(),
		Method: c.Method(),
		Time:   time.Now(),
		Reason: "Too many requests",
	}

	if _, err := os.Stat(rateLimitFile); err != nil {
		_ = os.MkdirAll("vulnerabilities", 0755)
		_, _ = os.Create(rateLimitFile)
	}

	file, err := os.OpenFile(rateLimitFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(logEntry)
}
