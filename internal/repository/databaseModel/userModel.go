package databaseModel

import (
	"Fiber_JWT_Authentication_backend_server/internal/controllers/serviceModels"
	"time"
)

type User struct {
	ID           int64     `db:"id"`
	Firstname    string    `db:"firstname"`
	Lastname     string    `db:"lastname"`
	Password     string    `db:"password"`
	Email        string    `db:"email"`
	Token        string    `db:"token"`
	UserType     string    `db:"user_type"`
	RefreshToken string    `db:"refresh_token"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (u *User) ToUserService() serviceModels.UserRequest {
	return serviceModels.UserRequest{
		ID:           u.ID,
		Firstname:    u.Firstname,
		Lastname:     u.Lastname,
		Password:     u.Password,
		Email:        u.Email,
		Token:        u.Token,
		UserType:     u.UserType,
		RefreshToken: u.RefreshToken,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
