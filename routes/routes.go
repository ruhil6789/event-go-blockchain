package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ruhil6789/event-sky/pkg/controllers"
)

func UserRoute(app *fiber.App) {
	//All routes related to users comes here
	app.Post("/user", controllers.CreateUser) //add this
	app.Get("/user/:userId", controllers.GetAUser) //get this route to get a single user data by id
}
