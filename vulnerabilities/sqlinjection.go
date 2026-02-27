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

var sqlRegex = regexp.MustCompile(`(?i)\b(union|select|insert|update|delete|drop|or|and|sleep|benchmark)\b|--|/\*|\*/|;|1=1|1=0`)

func IsSqlInjection(value string) bool {
	return sqlRegex.MatchString(value)
}

func SqlInjection(c *fiber.Ctx, query map[string]string, filename_sql string) bool {
	var vulner shared.Vulnerability
	found := false
	for i, j := range query {
		if IsSqlInjection(j) {
			found = true
			vulner.Key = i
			vulner.Value = j
			break
		}
	}
	if found {
		logSqlInjection(c, vulner, filename_sql)
	}

	return found
}

func SqlInjectionBody(c *fiber.Ctx, body map[string][]string, filename_sql string) bool {
	var vulner shared.Vulnerability

	found := false

	for key, values := range body {
		for _, value := range values {
			if IsSqlInjection(value) {
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
		logSqlInjection(c, vulner, filename_sql)
	}

	return found
}

func SqlInjectionInterface(c *fiber.Ctx, body map[string]interface{}, filename_sql string) bool {
	var vulner shared.Vulnerability
	found := false

	for key, value := range body {
		valStr := fmt.Sprintf("%v", value)
		if IsSqlInjection(valStr) {
			found = true
			vulner.Key = key
			vulner.Value = valStr
			break
		}
	}

	if found {
		logSqlInjection(c, vulner, filename_sql)
	}

	return found
}

func logSqlInjection(c *fiber.Ctx, vulner shared.Vulnerability, filename_sql string) {
	vulner.Ip = c.IP()
	vulner.Path = c.Path()
	vulner.Time = time.Now()

	if _, err := os.Stat(filename_sql); err != nil {
		os.Create(filename_sql)
	}

	file, err := os.OpenFile(filename_sql, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(vulner)
}
