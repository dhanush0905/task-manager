package routes

import "github.com/gofiber/fiber/v2"

func Setuproutes(app *fiber.App) {
	SetupTaskRoutes(app)
}
