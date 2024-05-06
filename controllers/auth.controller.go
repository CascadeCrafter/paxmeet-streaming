package controllers

import (
	"fmt"
	"streaming/initializers"
	"streaming/middleware"
	"streaming/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func GenerateToken(c *fiber.Ctx, config *initializers.Config) error {

	// Define the struct to get livekit Token
	type RequestData struct {
		RoomId     string `json:"roomId"`
		IsStreamer bool   `json:"isSreamer"`
	}

	user, ok := c.Locals("userDetails").(middleware.UserDetailsResponse)
	if !ok {
		// Handler case when user details are not properly set or wrong type
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Server error while retrieving user details",
		})
	}

	// Parse the POST request body into the struct
	var requestData RequestData
	if err := c.BodyParser(&requestData); err != nil {
		// Handle parsing error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Failed to parse request body: %v", err),
		})
	}

	livekitToken, err := utils.CreateToken(requestData.IsStreamer, requestData.RoomId, user.ID, user.Name, user.Photo, config)

	if err != nil {
		// Handler case when user details are not properly set or wrong type
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error while Generating Livekit Token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": livekitToken,
		},
	})
}

func RefreshToken(c *fiber.Ctx, config *initializers.Config) error {
	token := c.Get("token")
	tokenStr, err := jwt.ParseSigned(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to parse token",
		})
	}
	claims := &jwt.Claims{}
	secret := config.LiveKit.APISecret
	if err := tokenStr.Claims([]byte(secret), claims); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to parse token",
		})
	}
	expired := claims.Expiry.Time().Before(time.Now())
	if expired {
		return GenerateToken(c, config)
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"token": token,
			},
		})
	}
}
