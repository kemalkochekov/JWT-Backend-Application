package main

import (
	"Fiber_JWT_Authentication_backend_server/configs"
	"Fiber_JWT_Authentication_backend_server/internal/connectionDatabase"
	"Fiber_JWT_Authentication_backend_server/internal/connectionRedis"
	"Fiber_JWT_Authentication_backend_server/internal/repository/postgres"
	"Fiber_JWT_Authentication_backend_server/internal/repository/redis"
	"Fiber_JWT_Authentication_backend_server/internal/routes"
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Could not set up environment variable: %s", err)
	}

	httpPort := os.Getenv("PORT")

	dbConfig, err := configs.FromEnv()
	if err != nil {
		log.Fatalf("Could not import environment variables.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database, err := connectionDatabase.NewDB(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Could not connect database.")
	}

	defer func(database *connectionDatabase.Database) {
		err := database.Close()
		if err != nil {
			log.Printf("Error closing database: %s", err.Error())
		}
	}(database)

	if err != nil {
		log.Fatalf("Could not get environment variable: %v", err)
	}

	redisDatabase, err := connectionRedis.NewDatabase(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	redisClient := redis.NewClientRedisRepository(redisDatabase.Client)
	client := postgres.NewUserStorage(database)

	router := fiber.New()
	// Use the built-in Logger middleware
	router.Use(func(c *fiber.Ctx) error {
		// Log information about the incoming request
		println("Method:", c.Method(), "Path:", c.Path())
		return c.Next() // Move to the next middleware/handler
	})

	routes.AuthRoutes(router, &client, redisClient)
	routes.UserRoutes(router, &client, redisClient)

	err = router.Listen(httpPort)
	if err != nil {
		return
	}
}
