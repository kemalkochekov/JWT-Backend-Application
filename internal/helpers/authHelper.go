package helpers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func CheckUserType(ctx *fiber.Ctx, role string) error {
	userType := ctx.Locals("user_type")
	if userType != role {
		return errors.New("unauthorized to access this resource")
	}

	return nil
}
