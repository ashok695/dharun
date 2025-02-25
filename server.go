package main

import (
	"fmt"

	"github.com/dharun/poc/database"
	"github.com/dharun/poc/internals/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	database.DBConnection()
	defer database.CloseDatabase()
	app.Get("/planner", handlers.GetPlanner)
	fmt.Println(" 🚀APP is running on :3000")
	app.Listen(":3000")
}
