package middleware

import (
	"Fiber_JWT_Authentication_backend_server/internal/helpers"
	"Fiber_JWT_Authentication_backend_server/internal/repository/databaseModel"
	"Fiber_JWT_Authentication_backend_server/internal/utils"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func (m *MDWManager) Authenticate() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeaders := databaseModel.AuthHeaders{
			UserAgent:  ctx.Get(fiber.HeaderUserAgent),
			SessionKey: ctx.Cookies("token"),
		}
		// if empty, then clear cookie
		if authHeaders.SessionKey == "" {
			utils.ClearCookie(ctx, "token", "localhost")
			ctx.Status(fiber.StatusInternalServerError).Set("error", fmt.Sprintf("No Authorization header provided"))
			return errors.New("No Authorization header provided")
		}
		// if no session then clear cookie
		cachedSession, err := m.officiantRedisRepo.GetSession(context.Background(), authHeaders)
		if err != nil {
			utils.ClearCookie(ctx, "token", "localhost")
			ctx.Status(fiber.StatusInternalServerError).Set("error", fmt.Sprintf("No Authorization header provided"))
			return errors.New("No Authorization header provided")
		}

		claims, msg := helpers.ValidateToken(cachedSession.SessionKey)
		if msg != "" {
			ctx.Status(fiber.StatusInternalServerError).Set("error", msg)
			return errors.New(msg)
		}
		fmt.Println()
		// using Locals instead of set and get
		ctx.Locals("email", claims.Email)
		ctx.Locals("firstname", claims.Firstname)
		ctx.Locals("lastname", claims.Lastname)
		ctx.Locals("user_type", claims.UserType)
		log.Printf("Successfully Authenticated")
		return ctx.Next()
	}
}
