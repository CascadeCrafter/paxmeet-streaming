package controllers

import (
	"fmt"
	"streaming/initializers"
	"streaming/middleware"
	"streaming/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateTradingRoom(c *fiber.Ctx, config *initializers.Config) error {

	// Define the struct to get livekit Token
	type RequestData struct {
		RoomId   string `json:"roomId"`
		Products string `json:"products"`
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

	livekitToken, err := utils.CreateToken(true, requestData.RoomId, user.ID, user.Name, user.Photo, config)

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

func JoinTradingRoom(c *fiber.Ctx, config *initializers.Config) error {

	// Define the struct to get livekit Token
	type RequestData struct {
		RoomId string `json:"roomId"`
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

	livekitToken, err := utils.CreateToken(false, requestData.RoomId, user.ID, user.Name, user.Photo, config)

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

func GetTradingRoom(c *fiber.Ctx, config *initializers.Config) error {

}

func DeleteTradingRoom(c *fiber.Ctx, config *initializers.Config) error {

}
