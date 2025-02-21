package routes

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	SetupAuthRoutes(app)      // Authentication routes
	SetupTaskRoutes(app)      // Task management routes
	SetupWebSocketRoutes(app) // WebSocket for real-time updates
	SetupAIRoutes(app)        // AI-powered task breakdowns
}
