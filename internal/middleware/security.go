package middleware

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/neoway/golang-validation-documents/internal/config"
)

func SanitizeInput() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
			var body map[string]interface{}
			if err := c.BodyParser(&body); err == nil {
				sanitized := sanitizeMap(body)
				c.Locals("sanitized_body", sanitized)
			}
		}

		queryParams := make(map[string]string)
		c.Context().QueryArgs().VisitAll(func(key, value []byte) {
			sanitizedKey := sanitizeString(string(key))
			sanitizedValue := sanitizeString(string(value))
			queryParams[sanitizedKey] = sanitizedValue
		})
		c.Locals("sanitized_query", queryParams)

		return c.Next()
	}
}

func sanitizeMap(m map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})
	for k, v := range m {
		sanitizedKey := sanitizeString(k)
		switch val := v.(type) {
		case string:
			sanitized[sanitizedKey] = sanitizeString(val)
		case map[string]interface{}:
			sanitized[sanitizedKey] = sanitizeMap(val)
		default:
			sanitized[sanitizedKey] = val
		}
	}
	return sanitized
}

func sanitizeString(s string) string {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`[<>]`)
	s = re.ReplaceAllString(s, "")
	return s
}

func RequestCounter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		config.Metrics.IncrementRequestCount()
		return c.Next()
	}
}

