package routes

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App) {
	AuthRoute(app)
}