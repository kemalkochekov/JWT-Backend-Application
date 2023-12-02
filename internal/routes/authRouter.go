package routes

import (
	"Fiber_JWT_Authentication_backend_server/internal/connectionRedis"
	"Fiber_JWT_Authentication_backend_server/internal/controllers"
	"Fiber_JWT_Authentication_backend_server/internal/repository/postgres"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(incomingRoutes *fiber.App, client postgres.UserPgRepo, clientRedis connectionRedis.CacheRepository) {
	serviceUser := controllers.NewUserHandler(client, clientRedis)
	incomingRoutes.Post("users/signup", serviceUser.Signup())
	incomingRoutes.Post("users/login", serviceUser.Login())
}
