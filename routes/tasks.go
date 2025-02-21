package routes

import (
	"taskmanager/database"
	"taskmanager/models"
	"taskmanager/middleware"
	"github.com/gofiber/fiber/v2"
)

// Setup Task Routes
func SetupTaskRoutes(app *fiber.App) {
	taskGroup := app.Group("/tasks", middleware.JWTMiddleware)

	// ✅ Get All Tasks for Logged-in User
	taskGroup.Get("/", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)
		var tasks []models.Task
		database.DB.Where("user_id = ?", userID).Find(&tasks)
		return c.JSON(tasks)
	})

	// ✅ Create a Task (With WebSocket Update)
	taskGroup.Post("/", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)
		var task models.Task
		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		task.UserID = userID
		database.DB.Create(&task)

		// 🔥 Send WebSocket Update
		BroadcastUpdate("New task created: " + task.Title)

		return c.JSON(task)
	})

	// ✅ Update a Task (With WebSocket Update)
	taskGroup.Put("/:id", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)
		taskID := c.Params("id")

		var task models.Task
		if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}

		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		database.DB.Save(&task)

		// 🔥 Send WebSocket Update
		BroadcastUpdate("Task updated: " + task.Title)

		return c.JSON(task)
	})

	// ✅ Delete a Task (With WebSocket Update)
	taskGroup.Delete("/:id", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)
		taskID := c.Params("id")

		var task models.Task
		if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}

		database.DB.Delete(&task)

		// 🔥 Send WebSocket Update
		BroadcastUpdate("Task deleted: " + task.Title)

		return c.JSON(fiber.Map{"message": "Task deleted successfully"})
	})
}
