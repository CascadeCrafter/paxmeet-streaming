package api

import (
	"streaming/controllers"
	"streaming/initializers"
	"streaming/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(config *initializers.Config, micro *fiber.App) {
	micro.Route("/livekit", func(router fiber.Router) {
		router.All("/*", middleware.CheckAuth(config.Auth.Uri), func(c *fiber.Ctx) error {
			return controllers.LivekitHandler(c, config)
		})
	})

	micro.Route("/auth", func(router fiber.Router) {
		router.Post("/token", middleware.CheckAuth(config.Auth.Uri), func(c *fiber.Ctx) error {
			return controllers.GenerateToken(c, config)
		})
	})

	// Health check route
	micro.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Service is up and running!",
		})
	})

	micro.Route("/streaming", func(router fiber.Router) {
		router.Post("/room/create", middleware.CheckAuth(config.Auth.Uri), func(c *fiber.Ctx) error {
			return controllers.CreateTradingRoom(c, config)
		})
		router.Post("/room/entry", middleware.CheckAuth(config.Auth.Uri), func(c *fiber.Ctx) error {
			return controllers.EntryTradingRoom(c, config)
		})
		router.Get("/room/get/:roomId",
			// middleware.CheckAuth(config.Auth.Uri),
			func(c *fiber.Ctx) error {
				return controllers.GetTradingRoom(c, config)
			})
		router.Get("/rooms/all",
			// middleware.CheckAuth(config.Auth.Uri),
			func(c *fiber.Ctx) error {
				return controllers.GetAllTradingRooms(c, config)
			})
		router.Post("/room/join/:roomId",
			// middleware.CheckAuth(config.Auth.Uri),
			func(c *fiber.Ctx) error {
				return controllers.JoinTradingRoom(c, config)
			})
		router.Delete("/room/delete/:roomId", middleware.CheckAuth(config.Auth.Uri), func(c *fiber.Ctx) error {
			return controllers.DeleteTradingRoom(c, config)
		})
		router.Post("/checkTokenExp", middleware.CheckAuth(config.Auth.Uri), func(c *fiber.Ctx) error {
			return controllers.RefreshToken(c, config)
		})
		router.Get("/rooms/get", func(c *fiber.Ctx) error {
			return controllers.GetRooms(c, config)
		})
	})
}
