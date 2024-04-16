package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CheckAuth(auth_uri string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}
		req, err := http.NewRequest("GET", auth_uri, nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
		}

		req.Header.Add("Authorization", token)
		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return c.Status(fiber.StatusForbidden).SendString("Forbidden")
		}

		return c.Next()
	}
}
