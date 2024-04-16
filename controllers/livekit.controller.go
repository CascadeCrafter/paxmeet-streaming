package controllers

import (
	"fmt"
	"streaming/initializers"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func LivekitHandler(c *fiber.Ctx, config *initializers.Config) error {

	// Construct the target URL by appending the original request path.
	baseURL := strings.TrimRight(config.LiveKit.Uri, "/")
	targetPath := strings.TrimLeft(c.Params("*"), "/") // Params("*") captures all after /livekit
	targetURL := fmt.Sprintf("%s/%s", baseURL, targetPath)

	return proxy.Do(c, targetURL)
}
