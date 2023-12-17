package routes

import (
	"Fiber_JWT_Authentication_backend_server/internal/controllers"
	"Fiber_JWT_Authentication_backend_server/internal/repository/postgres"
	"Fiber_JWT_Authentication_backend_server/pkg/connectionRedis"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(incomingRoutes *fiber.App, client postgres.UserPgRepo, clientRedis connectionRedis.CacheRepository) {
	serviceUser := controllers.NewUserHandler(client, clientRedis)
	incomingRoutes.Get("/", serviceUser.Main())
	incomingRoutes.Post("users/signup", serviceUser.Signup())
	incomingRoutes.Post("users/login", serviceUser.Login())
}
