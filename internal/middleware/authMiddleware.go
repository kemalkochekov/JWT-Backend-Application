package middleware

import (
	"Fiber_JWT_Authentication_backend_server/internal/helpers"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func Authenticate() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// getting token from header
		clientToken := ctx.Get("token")
		if clientToken == "" {
			ctx.Status(fiber.StatusInternalServerError).Set("error", fmt.Sprintf("No Authorization header provided"))
			return errors.New("No Authorization header provided")
		}
		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			ctx.Status(fiber.StatusInternalServerError).Set("error", err)
			return errors.New(err)
		}
		// using Locals instead of set and get
		ctx.Locals("email", claims.Email)
		ctx.Locals("firstname", claims.Firstname)
		ctx.Locals("lastname", claims.Lastname)
		ctx.Locals("user_type", claims.UserType)
		log.Printf("Successfully Authenticated")
		return ctx.Next()
	}
}
