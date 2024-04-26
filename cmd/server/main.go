package main

import (
	"fmt"
	"log"
	"os"
	"streaming/api"
	"streaming/initializers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// main function
func main() {
	config, err := initializers.LoadConfig()
	initializers.ConnectRedis(config)
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("config: ", *config)

	app := fiber.New(fiber.Config{
		ServerHeader: "PaxStreaming",
		BodyLimit:    20 * 1024 * 1024, // 20 MB
	})

	micro := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Access-Control-Allow-Headers, Session, Mode",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: false,
	}))

	api.Register(config, micro)

	app.Mount("/", micro)

	log.Fatal(app.Listen(":8080"))
}
