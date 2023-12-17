package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mysql/handlers"
)

func SetUpRoute(app *fiber.App) {
	employee := app.Group("/api")
	employee.Post("/create", handlers.Create)
	employee.Get("/employees", handlers.GetAll)
	employee.Get("/employee/:id", handlers.GetById)
	employee.Put("/employee/:id", handlers.Update)
	employee.Delete("/employee/:id", handlers.Delete)
	employee.Get("/image/:filename", handlers.GetImage)
}
