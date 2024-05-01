package controllers

import (
	"encoding/json"
	"fmt"
	"streaming/initializers"
	"streaming/middleware"
	"streaming/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// Creating a struct to encapsulate both RoomId, Products, and PublisherId
type RoomDetails struct {
	Products    []string `json:"products"`
	PublisherId string   `json:"publisherId"`
	Title       string   `json:"title"`
}

func CreateTradingRoom(c *fiber.Ctx, config *initializers.Config) error {

	// Define the struct to get livekit Token
	type RequestData struct {
		RoomId   string   `json:"roomId"`
		Products []string `json:"products"`
		Title    string   `json:"title"`
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

	// Assign values to RoomDetails
	roomDetails := RoomDetails{
		Products:    requestData.Products,
		PublisherId: user.ID,
		Title:       requestData.Title,
	}

	// Convert roomDetails into a JSON string
	roomDetailsJSON, err := json.Marshal(roomDetails)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to serialize room details",
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

	// After successfully generating a livekit token and before returning success response:
	err = initializers.RedisClient.Set(initializers.Ctx, "room:"+requestData.RoomId, roomDetailsJSON, 12*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to store room data in Redis",
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
	roomId := c.Params("roomId")
	// Fetch room details from Redis
	val, err := initializers.RedisClient.Get(initializers.Ctx, "room:"+roomId).Result()

	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Room not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error fetching room details",
		})
	}

	// Assuming val is the JSON string, deserialize it
	var roomDetails RoomDetails
	err = json.Unmarshal([]byte(val), &roomDetails)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to deserialize room details",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   roomDetails,
	})

}

func GetAllTradingRooms(c *fiber.Ctx, config *initializers.Config) error {
	// WARNING: Use of KEYS in production environments with large data sets is discouraged.
	keys, err := initializers.RedisClient.Keys(initializers.Ctx, "room:*").Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error fetching rooms",
		})
	}

	// Prepare a map to hold room data
	rooms := make(map[string]RoomDetails)
	for _, key := range keys {
		val, err := initializers.RedisClient.Get(initializers.Ctx, key).Result()
		if err != nil {
			continue // Optionally log this error
		}

		var roomDetails RoomDetails
		if err := json.Unmarshal([]byte(val), &roomDetails); err != nil {
			continue // Optionally log this error
		}

		// Extract roomId from the key and use it as a map key
		roomId := strings.TrimPrefix(key, "room:")
		rooms[roomId] = roomDetails
	}

	// Encode the entire map as a JSON object
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"rooms":  rooms,
	})
}

func DeleteTradingRoom(c *fiber.Ctx, config *initializers.Config) error {
	roomId := c.Params("roomId")

	err := initializers.RedisClient.Del(initializers.Ctx, "room:"+roomId).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Failed to delete room data from Redis: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}
