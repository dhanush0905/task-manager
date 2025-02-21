package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"taskmanager/database" // Replace 'yourmodule' with your updated module name
	"taskmanager/models"
)

// Get all tasks
func GetAllTasks(c *fiber.Ctx) error {
	var tasks []models.Task
	database.DB.Find(&tasks)
	return c.JSON(tasks)
}

// Create a new task
func CreateTask(c *fiber.Ctx) error {
	task := new(models.Task)
	if err := c.BodyParser(task); err != nil {
		fmt.Println("❌ Error parsing request body:", err) // ✅ Debugging log
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	database.DB.Create(&task)
	return c.Status(201).JSON(task)
}

// Update an existing task
func UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	database.DB.Save(&task)
	return c.JSON(task)
}

// Delete a task
func DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}
	database.DB.Delete(&task)
	return c.JSON(fiber.Map{"message": "Task deleted"})
}
