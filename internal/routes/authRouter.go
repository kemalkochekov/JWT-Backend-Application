package routes

import (
	"Fiber_JWT_Authentication_backend_server/internal/controllers"
	"Fiber_JWT_Authentication_backend_server/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(incomingRoutes *fiber.App, client repository.UserPgRepo) {
	serviceUser := controllers.NewUserHandler(client)
	incomingRoutes.Post("users/signup", serviceUser.Signup())
	incomingRoutes.Post("users/login", serviceUser.Login())
}
