package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CheckAuth(auth_uri string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			fmt.Println("Unauthorized")
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

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error reading response body from auth server"})
		}
		defer resp.Body.Close()

		var userDetails map[string]interface{}
		if err := json.Unmarshal(body, &userDetails); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error parsing response body"})
		}

		// Store user details in Locals for access in subsequent handlers
		c.Locals("user", userDetails)

		fmt.Println("Authentication Passed...")

		return c.Next()
	}
}
