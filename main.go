package main

import (
	"style-stamp/config"
	"style-stamp/pkg/routes"

	"github.com/gofiber/fiber/v2"
)


func main(){
	config.ConnectDb()
	app := fiber.New()
	routes.Routes(app)
	app.Get("/",func (c *fiber.Ctx)error  {
		return c.JSON(fiber.Map{
			"status":200,
			"message": "Hello, World!",
		})
	})
	app.Listen(":5050")
}