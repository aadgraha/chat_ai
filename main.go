package main

import (
	"chat_ai/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	routes.SetupRoutes(app)
	port := os.Getenv("PORT")
	if err := app.Listen(port); err != nil {
		panic(err)
	}
}
