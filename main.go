package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ruhil6789/event-sky/database"
	"github.com/ruhil6789/event-sky/routes"
)

func main() {
	fmt.Println("Hello, Go!")

	app := fiber.New()

	database.ConnectMongoDb()
	routes.UserRoute(app)
	app.Listen(":8080")
}
