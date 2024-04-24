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
}
