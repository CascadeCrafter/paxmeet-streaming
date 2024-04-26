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

	micro.Route("/streaming", func(router fiber.Router) {
		micro.Post("/room/create", func(c *fiber.Ctx) error {
			return controllers.CreateTradingRoom(c, config)
		})
		micro.Get("/room/get", func(c *fiber.Ctx) error {
			return controllers.GetTradingRoom(c, config)
		})

		micro.Post("/room/join", func(c *fiber.Ctx) error {
			return controllers.JoinTradingRoom(c, config)
		})

		micro.Delete("/room/delete", func(c *fiber.Ctx) error {
			return controllers.DeleteTradingRoom(c, config)
		})
	})
}
