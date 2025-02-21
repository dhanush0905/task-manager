package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // for loading .env
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"taskmanager/database"
	"taskmanager/routes"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// Set port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Verify JWT_SECRET is loaded
	fmt.Println("✅ JWT_SECRET Loaded Successfully:", os.Getenv("JWT_SECRET"))

	// Initialize Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Connect to Database
	database.ConnectDB()

	// Default route for testing
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Task Manager API is running successfully!")
	})

	// Setup API Routes (Authentication & Task routes)
	routes.SetupRoutes(app)

	// Start Server
	log.Printf("✅ Server running at http://127.0.0.1:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
