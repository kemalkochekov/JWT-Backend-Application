package routes

import (
	"Fiber_JWT_Authentication_backend_server/internal/connectionRedis"
	"Fiber_JWT_Authentication_backend_server/internal/controllers"
	"Fiber_JWT_Authentication_backend_server/internal/middleware"
	"Fiber_JWT_Authentication_backend_server/internal/repository/postgres"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(incomingRoutes *fiber.App, client postgres.UserPgRepo, clientRedis connectionRedis.CacheRepository) {
	serviceUser := controllers.NewUserHandler(client, clientRedis)
	mw := middleware.NewOfficiantMiddleware(clientRedis)
	incomingRoutes.Use(mw.Authenticate())
	incomingRoutes.Get("/admin", serviceUser.GetUsers())
	incomingRoutes.Get("/users", serviceUser.GetUser())
	incomingRoutes.Get("/logout", serviceUser.Logout())
}
