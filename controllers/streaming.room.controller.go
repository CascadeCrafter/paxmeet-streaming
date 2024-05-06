package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	Products  json.RawMessage                `json:"products"`
	Publisher middleware.UserDetailsResponse `json:"publisher"`
	Title     string                         `json:"title"`
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
	requestData := new(RequestData)
	if err := c.BodyParser(requestData); err != nil {
		// Handle parsing error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Failed to parse request body: %v", err),
		})
	}
	
	// ** Here Fetch data from backend with requestData.Products and store Products **
	// Use the new function to fetch product details
	fetchedProducts, err := fetchProductDetailsFromBackend(requestData.Products, user.ID)
	if err != nil {
		// Handle error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Error fetching product details: %v", err),
		})
	}

	// Assign values to RoomDetails
	roomDetails := RoomDetails{
		Products:  fetchedProducts,
		Publisher: user,
		Title:     requestData.Title,
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

	// // Define the struct to get livekit Token
	
	type RequestData struct {
		UserId string `json:userId`
		UserPhoto string `json:"photo"`
		UserName string `json:"userName"`
	}


	roomId := c.Params("roomId")

	// Fetch room details from Redis
	// _, err := initializers.RedisClient.Get(initializers.Ctx, "room:"+roomId).Result()

	// if err == redis.Nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Room not found",
	// 	})
	// }
	// user, ok := c.Locals("userDetails").(middleware.UserDetailsResponse)
	// if !ok {
	// 	// Handler case when user details are not properly set or wrong type
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Server error while retrieving user details",
	// 	})
	// }

	// // Parse the POST request body into the struct
	requestData:= &RequestData{}
	// var userId UserId
	
	if err := c.BodyParser(requestData); err != nil {
		// Handle parsing error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Failed to parse request body: %v", err),
		})
	}

	livekitToken, err := utils.CreateToken(false, roomId, requestData.UserId, requestData.UserName, requestData.UserPhoto, config)

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

func fetchProductDetailsFromBackend(productIDs []string, publisherID string) ([]byte, error) {
	// Construct the request payload
	requestPayload := struct {
		IDs       []string `json:"ids"`
		Publisher string   `json:"publisher"`
	}{
		IDs:       productIDs,
		Publisher: publisherID,
	}

	// Convert the request payload to JSON bytes
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Execute the POST request to the backend service
	resp, err := http.Post("https://go.paxintrade.com/api/blog/filterByIds", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from backend: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is not 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend responded with non-200 status code: %d", resp.StatusCode)
	}

	// Assuming the structure of the response for our case is known and consistent
	// Decode the entire response to access the Blogs part
	var backendResponse struct {
		Blogs  json.RawMessage `json:"blogs"` // Use json.RawMessage to get the raw JSON
		Status string          `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&backendResponse); err != nil {
		return nil, fmt.Errorf("failed to decode backend response: %w", err)
	}

	// If needed, perform further processing on backendResponse.Blogs, which is raw JSON

	// Return the raw JSON of the Blogs part directly
	return backendResponse.Blogs, nil
}
