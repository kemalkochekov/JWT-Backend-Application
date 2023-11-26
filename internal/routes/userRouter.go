package routes

import (
	"Fiber_JWT_Authentication_backend_server/internal/controllers"
	"Fiber_JWT_Authentication_backend_server/internal/middleware"
	"Fiber_JWT_Authentication_backend_server/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(incomingRoutes *fiber.App, client repository.UserPgRepo) {
	serviceUser := controllers.NewUserHandler(client)
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.Get("/users", serviceUser.GetUsers())
	incomingRoutes.Get("/users/:user_id", serviceUser.GetUser())
}
