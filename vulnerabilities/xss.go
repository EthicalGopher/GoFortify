package vulnerabilities

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/EthicalGopher/GoFortify/shared"
	"github.com/gofiber/fiber/v2"
)

var xssRegex = regexp.MustCompile(`(?i)<script|</script>|javascript:|onerror=|onload=|onclick=|onmouseover=|onfocus=|<img|<svg|<iframe|<body|<html|<meta|alert\(|confirm\(|prompt\(|document\.cookie|document\.location|window\.location|eval\(|base64,`)

func IsXSS(value string) bool {
	return xssRegex.MatchString(value)
}

func XSSInjection(c *fiber.Ctx, query map[string]string, filename_xss string) bool {
	var vulner shared.Vulnerability
	found := false
	for key, value := range query {
		if IsXSS(value) {
			found = true
			vulner.Key = key
			vulner.Value = value
			break
		}
	}
	if found {
		logXSS(c, vulner, filename_xss)
	}
	return found
}

func XSSInjectionBody(c *fiber.Ctx, body map[string][]string, filename_xss string) bool {
	var vulner shared.Vulnerability
	found := false
	for key, values := range body {
		for _, value := range values {
			if IsXSS(value) {
				found = true
				vulner.Key = key
				vulner.Value = value
				break
			}
		}
		if found {
			break
		}
	}
	if found {
		logXSS(c, vulner, filename_xss)
	}
	return found
}

func XSSInjectionInterface(c *fiber.Ctx, body map[string]interface{}, filename_xss string) bool {
	var vulner shared.Vulnerability
	found := false
	for key, value := range body {
		valStr := fmt.Sprintf("%v", value)
		if IsXSS(valStr) {
			found = true
			vulner.Key = key
			vulner.Value = valStr
			break
		}
	}
	if found {
		logXSS(c, vulner, filename_xss)
	}
	return found
}

func logXSS(c *fiber.Ctx, vulner shared.Vulnerability, filename_xss string) {
	vulner.Ip = c.IP()
	vulner.Path = c.Path()
	vulner.Time = time.Now()

	if _, err := os.Stat(filename_xss); err != nil {
		os.Create(filename_xss)
	}

	file, err := os.OpenFile(filename_xss, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(vulner)
	}
}
