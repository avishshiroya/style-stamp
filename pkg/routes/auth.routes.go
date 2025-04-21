package routes

import (
	"style-stamp/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App){
	route := app.Group("/auth")
	route.Post("/",controllers.SignUp)
}