package main

import (
	"Fiber_JWT_Authentication_backend_server/configs"
	"Fiber_JWT_Authentication_backend_server/internal/repository/postgres"
	"Fiber_JWT_Authentication_backend_server/internal/repository/redis"
	"Fiber_JWT_Authentication_backend_server/internal/routes"
	"Fiber_JWT_Authentication_backend_server/pkg/connectionDatabase"
	"Fiber_JWT_Authentication_backend_server/pkg/connectionRedis"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

const contextTimeOut = 10 * time.Second

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Could not set up environment variable: %s", err)
	}

	httpPort := os.Getenv("PORT")

	dbConfig, err := configs.FromEnv()
	if err != nil {
		log.Fatalf("Could not import environment variables.")
	}

	redisConfig, err := configs.FromEnvRedis()
	if err != nil {
		log.Fatalf("Could not import environment variables for Redis.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeOut)
	defer cancel()

	database, err := connectionDatabase.NewDB(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Could not connect database because of %v.", err)
	}

	defer func(database *connectionDatabase.Database) {
		err := database.Close()
		if err != nil {
			log.Printf("Error closing Postgres database: %s", err.Error())
		}
	}(database)

	if err != nil {
		log.Fatalf("Could not get environment variable: %v", err)
	}

	redisDatabase, err := connectionRedis.NewDatabase(ctx, redisConfig)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	defer func(redisDatabase *connectionRedis.Database) {
		err := redisDatabase.Client.Close()
		if err != nil {
			log.Printf("Error closing Redis database: %s", err.Error())
		}
	}(redisDatabase)

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(quit)

	go func() {
		<-quit
		// Perform graceful shut down

		err := database.Close()
		if err != nil {
			log.Printf("Error closing database: %s", err.Error())
		}

		err = redisDatabase.Client.Close()
		if err != nil {
			log.Printf("Error closing redis database: %s", err.Error())
		}

		os.Exit(0)
	}()

	err = router.Listen(httpPort)
	if err != nil {
		return
	}
}
