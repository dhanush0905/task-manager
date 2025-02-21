package routes

import (
	"fmt"
	"os"
	"time"

	"taskmanager/database"
	"taskmanager/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// SetupAuthRoutes registers authentication routes
func SetupAuthRoutes(app *fiber.App) {
	authGroup := app.Group("/auth")
	authGroup.Post("/signup", Signup)
	authGroup.Post("/login", Login)
}

// Signup registers a new user
func Signup(c *fiber.Ctx) error {
	user := new(models.User)

	// Parse JSON body
	if err := c.BodyParser(user); err != nil {
		fmt.Println("❌ Error parsing request:", err) // Debugging log
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check if email or username already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).Or("username = ?", user.Username).First(&existingUser).Error; err == nil {
		fmt.Println("❌ Error: Email or username already taken") // Debugging log
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email or username already in use"})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("❌ Error hashing password:", err) // Debugging log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	user.Password = string(hashedPassword)

	// Save user to the database
	result := database.DB.Create(&user)
	if result.Error != nil {
		fmt.Println("❌ Database Error:", result.Error) // Debugging log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	fmt.Println("✅ User registered successfully:", user.Email) // Debugging log
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

// Login authenticates a user and returns a JWT token
func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON body
	if err := c.BodyParser(&input); err != nil {
		fmt.Println("❌ Error parsing login request:", err) // Debugging log
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("❌ Error: Email not found") // Debugging log
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
		}
		fmt.Println("❌ Database Error:", err) // Debugging log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		fmt.Println("❌ Error: Incorrect password") // Debugging log
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Generate JWT token
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		fmt.Println("❌ Error: JWT_SECRET is not set!") // Debugging log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error, please try again later"})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Println("❌ Error generating token:", err) // Debugging log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	fmt.Println("✅ Token generated successfully for:", user.Email) // Debugging log
	return c.JSON(fiber.Map{"token": tokenString})
}
