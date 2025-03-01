package main

import (
	"github.com/dharun/poc/database"
	"github.com/dharun/poc/internals/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	// Initialize Fiber app
	app := fiber.New()
	app.Use(cors.New())

	// Initialize Database
	database.DBConnection()
	defer database.CloseDatabase()
	// Define Fiber routes
	app.Get("/planner", handlers.GetPlanner)
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Fiber server failed: %v", err)
	}
}
