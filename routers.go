package main

import "github.com/gofiber/fiber/v2"

func InitRoutes(app *fiber.App) {
	app.Get("/converter/:from_currency", syncConverter)
	app.Get("/converter/async/:from_currency", asyncConverter)
	app.Get("/converter/async/v2/:from_currency", asyncConverterV2)
}
