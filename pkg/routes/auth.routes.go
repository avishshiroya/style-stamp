package routes

import (
	"style-stamp/app/controllers"
	"style-stamp/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App) {
	route := app.Group("/auth")
	route.Post("/", controllers.SignUp)
	route.Post("/login", controllers.SignIn)
	route.Get("/", middlewares.Authentication(), controllers.GetDetails)
	route.Get("/logout", middlewares.Authentication(), controllers.LogOut)
	route.Get("/refresh-token", controllers.RefreshToken)
}
