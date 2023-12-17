package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mysql/database"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mysql/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed To Load Env")
	}

	database.ConnectDB()
	app := fiber.New()
	routes.SetUpRoute(app)

	app.Listen(":8080")
}
