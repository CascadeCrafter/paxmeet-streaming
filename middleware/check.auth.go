package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Define userDetailsResponse according to the structure we expect to receive from the auth server
type UserDetailsResponse struct {
	// ID           uuid.UUID `json:"userID"`
	ID           string `json:"userID"`
	Photo        string `json:"photo"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	TelegramName string `json:"telegramname"`
}

type APIResponse struct {
	Status string              `json:"status"`
	Data   UserDetailsResponse `json:"data"`
}

func CheckAuth(authURI string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		req, err := http.NewRequest("GET", authURI, nil)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
		}
		req.Header.Add("Authorization", token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Forbidden"})
		}

		defer resp.Body.Close()

		// If status is not OK, read the body and forward the auth service's response
		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error reading auth service response"})
			}

			// Parse the response to check if it's a valid JSON
			var respContent interface{}
			if err := json.Unmarshal(body, &respContent); err != nil {
				// If not a JSON or failed to parse, forward as plain text
				return c.Status(resp.StatusCode).SendString(string(body))
			}

			// If the response is JSON, forward it as such
			return c.Status(resp.StatusCode).JSON(respContent)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error reading response body"})
		}

		// Unmarshal the body into the userDetailsResponse struct
		var apiResponse APIResponse

		if err := json.Unmarshal(body, &apiResponse); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error parsing response body"})
		}

		// Store the structured userDetails into Locals
		c.Locals("userDetails", apiResponse.Data)

		fmt.Println("Authentication Passed...")
		return c.Next()
	}
}
