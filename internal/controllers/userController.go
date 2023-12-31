package controllers

import (
	"Fiber_JWT_Authentication_backend_server/internal/controllers/serviceModels"
	"Fiber_JWT_Authentication_backend_server/internal/helpers"
	"Fiber_JWT_Authentication_backend_server/internal/repository/databaseModel"
	"Fiber_JWT_Authentication_backend_server/internal/repository/postgres"
	"Fiber_JWT_Authentication_backend_server/pkg/connectionRedis"
	"Fiber_JWT_Authentication_backend_server/pkg/utils"
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const (
	generatePasswordCost = 14
	contextTimeOut       = 50 * time.Second
	cookieTimeOut        = 15 * time.Minute
	sessionTimeOut       = 15
)

var validate = validator.New()

type UserHandler struct {
	userRepo      postgres.UserPgRepo
	userRedisRepo connectionRedis.CacheRepository
}

func NewUserHandler(clientRepo postgres.UserPgRepo, clientRedisRepo connectionRedis.CacheRepository) *UserHandler {
	return &UserHandler{userRepo: clientRepo, userRedisRepo: clientRedisRepo}
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), generatePasswordCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return false, "Email or Password is incorrect!!!"
	}

	return true, ""
}

func (u *UserHandler) Main() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("Welcome JWT Backend Application")
	}
}

func (u *UserHandler) Signup() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var user serviceModels.UserRequest

		err := ctx.BodyParser(&user)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest).Set("error", err.Error())
			return err
		}

		validationErr := validate.Struct(&user)
		if validationErr != nil {
			ctx.Status(fiber.StatusBadRequest).Set("error", validationErr.Error())
			return validationErr
		}

		user.CreatedAt, err = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).Set("error", err.Error())
			return err
		}

		user.UpdatedAt, err = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).Set("error", err.Error())
			return err
		}

		user.Token, user.RefreshToken, err = helpers.GenerateAllTokens(user.Email, user.Firstname, user.Lastname, user.ID, user.UserType)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).Set("error", err.Error())
			return err
		}

		user.Password, err = HashPassword(user.Password)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).Set("error", err.Error())
			return err
		}

		ctxDb, cancel := context.WithTimeout(context.Background(), contextTimeOut)

		defer cancel()

		err = u.userRepo.RegisterUser(ctxDb, user)
		if err != nil {
			if errors.Is(err, errors.New("Email already exists")) {
				ctx.Status(fiber.StatusBadRequest).Set("error", err.Error())
				return err
			}

			ctx.Status(fiber.StatusInternalServerError).Set("error", err.Error())

			return err
		}

		return ctx.Status(fiber.StatusOK).SendString("Successfully Registered User Account")
	}
}

func (u *UserHandler) Login() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var user serviceModels.UserRequest

		err := ctx.BodyParser(&user)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest).Set("error", err.Error())
			return err
		}

		providedUserPassword := user.Password
		ctxDb, cancel := context.WithTimeout(context.Background(), contextTimeOut)

		defer cancel()

		user, err = u.userRepo.GetUser(ctxDb, user.Email)
		if err != nil {
			ctx.Status(fiber.StatusNotFound).Set("error", err.Error())
			return err
		}

		token, refreshToken, err := helpers.GenerateAllTokens(user.Email, user.Firstname, user.Lastname, user.ID, user.UserType)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).Set("error", err.Error())
			return err
		}

		err = u.userRepo.LoginUser(ctxDb, user.Email, token, refreshToken)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).Set("error", err.Error())
			return err
		}

		user.Token = token
		user.RefreshToken = refreshToken

		passwordIsValid, msg := VerifyPassword(user.Password, providedUserPassword)
		if !passwordIsValid {
			ctx.Status(fiber.StatusInternalServerError).Set("error", msg)
			return errors.New(msg)
		}

		userAgent := ctx.Get("User-Agent")

		utils.SetCookie(ctx, utils.CookieData{
			Name:    "token",
			Value:   refreshToken,
			Expires: time.Now().Add(cookieTimeOut),
			Domain:  "localhost",
		})

		err = u.userRedisRepo.PutSession(ctxDb, databaseModel.CacheUserSession{
			UserAgent:  userAgent,
			SessionKey: refreshToken,
			Duration:   time.Duration(sessionTimeOut) * time.Minute,
			CreatedAt:  time.Now().UnixMilli(),
		})

		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).Set("error", err.Error())
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(user)
	}
}
func (u *UserHandler) GetUsers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Locals("email") == nil || ctx.Locals("firstname") == nil || ctx.Locals("lastname") == nil || ctx.Locals("user_type") == nil {
			ctx.Status(fiber.StatusBadRequest).Set("error", "Unauthorized to access this resource")
			return errors.New("Unauthorized to access this resource")
		}

		if err := helpers.CheckUserType(ctx, "ADMIN"); err != nil {
			ctx.Status(fiber.StatusBadRequest).Set("error", err.Error())
			return errors.New("Only Admin has access for this resource")
		}

		ctxDb, cancel := context.WithTimeout(context.Background(), contextTimeOut)

		defer cancel()

		users, err := u.userRepo.AdminGetUsers(ctxDb)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).Set("error", err.Error())
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(users)
	}
}

func (u *UserHandler) GetUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Locals("email") == nil || ctx.Locals("firstname") == nil || ctx.Locals("lastname") == nil || ctx.Locals("user_type") == nil {
			ctx.Status(fiber.StatusBadRequest).Set("error", "Unauthorized to access this resource")
			return errors.New("Unauthorized to access this resource")
		}

		userEmail, ok := ctx.Locals("email").(string)
		if !ok {
			return errors.New("could not convert into string")
		}
		ctxDb, cancel := context.WithTimeout(context.Background(), contextTimeOut)

		defer cancel()

		var user serviceModels.UserRequest

		user, err := u.userRepo.GetUser(ctxDb, userEmail)
		if err != nil {
			ctx.Status(fiber.StatusNotFound).Set("error", err.Error())
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(user)
	}
}

func (u *UserHandler) Logout() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sessionKey := ctx.Cookies("token")
		if sessionKey == "" {
			ctx.Status(fiber.StatusBadRequest).Set("error", "Unauthorized user")
			return errors.New("Unauthorized user")
		}

		utils.ClearCookie(ctx, "token", "127.0.0.1")

		return ctx.Status(fiber.StatusOK).SendString("Successfully Logged out")
	}
}
