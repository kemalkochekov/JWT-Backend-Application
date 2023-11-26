package main

import (
	"Fiber_JWT_Authentication_backend_server/internal/configs"
	"Fiber_JWT_Authentication_backend_server/internal/repository"
	"Fiber_JWT_Authentication_backend_server/internal/repository/connectionDatabase"
	"Fiber_JWT_Authentication_backend_server/internal/routes"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Could not set up environment variable: %s", err)
	}
	httpPort := os.Getenv("PORT")
	dbConfig, err := configs.FromEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database, err := connectionDatabase.NewDB(ctx, dbConfig)
	defer database.Close()
	if err != nil {
		log.Fatalf("Could not get environment variable: %v", err)
	}
	client := repository.NewUserStorage(database)
	router := fiber.New()
	// Use the built-in Logger middleware
	router.Use(func(c *fiber.Ctx) error {
		// Log information about the incoming request
		println("Method:", c.Method(), "Path:", c.Path())
		return c.Next() // Move to the next middleware/handler
	})

	routes.AuthRoutes(router, &client)
	routes.UserRoutes(router, &client)

	router.Get("/api-1", func(ctx *fiber.Ctx) error {
		ctx.Set("success", "Access granted for api-1")
		return ctx.Status(fiber.StatusOK).SendString("")
	})
	router.Get("/api-2", func(ctx *fiber.Ctx) error {
		ctx.Set("success", "Access granted for api-2")
		return ctx.Status(fiber.StatusOK).SendString("")
	})
	router.Listen(httpPort)
}
