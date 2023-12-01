package postgres

import (
	"Fiber_JWT_Authentication_backend_server/internal/connectionDatabase"
	"Fiber_JWT_Authentication_backend_server/internal/controllers/serviceModels"
	"Fiber_JWT_Authentication_backend_server/internal/repository/databaseModel"
	"Fiber_JWT_Authentication_backend_server/internal/utils"
	"context"
	"errors"
	"time"
)

type UserPgRepo interface {
	RegisterUser(ctx context.Context, userReq serviceModels.UserRequest) error
	LoginUser(ctx context.Context, email string, token string, refreshToken string) error
	GetUser(ctx context.Context, email string) (serviceModels.UserRequest, error)
	AdminGetUsers(ctx context.Context) ([]serviceModels.UserRequest, error)
}
type UserStorage struct {
	db connectionDatabase.DBops
}

func NewUserStorage(database connectionDatabase.DBops) UserStorage {
	return UserStorage{db: database}
}

func ToOfficiantStorage(s serviceModels.UserRequest) databaseModel.User {
	return databaseModel.User{
		Firstname:    s.Firstname,
		Lastname:     s.Lastname,
		Password:     s.Password,
		Email:        s.Email,
		Token:        s.Token,
		UserType:     s.UserType,
		RefreshToken: s.RefreshToken,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

//service into database convert func if needed

func (d *UserStorage) RegisterUser(ctx context.Context, userReq serviceModels.UserRequest) error {
	user := ToOfficiantStorage(userReq)
	var exists bool
	err := d.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Email or Phone number already exists")
	}
	_, err = d.db.ExecContext(ctx, `INSERT INTO users(firstname, lastname, password, email, token, user_type, refresh_token, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);`,
		user.Firstname,
		user.Lastname,
		user.Password,
		user.Email,
		user.Token,
		user.UserType,
		user.RefreshToken,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *UserStorage) LoginUser(ctx context.Context, email string, token string, refreshToken string) error {
	updated_at, _ := time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
	_, err := d.db.ExecContext(ctx,
		`UPDATE users SET token = $2, refresh_token = $3, updated_at = $4 WHERE email = $1`, email, token, refreshToken, updated_at)
	if err != nil {
		return err
	}
	return nil
}
func (d *UserStorage) AdminGetUsers(ctx context.Context) ([]serviceModels.UserRequest, error) {
	var users []databaseModel.User
	err := d.db.SelectContext(ctx, &users, `SELECT id, firstname, lastname, password, email, user_type, created_at, updated_at, token, refresh_token FROM users`)
	if err != nil {
		return nil, err
	}
	usersService := utils.Map(
		users,
		func(item databaseModel.User) serviceModels.UserRequest {
			return item.ToUserService()
		},
	)
	return usersService, nil
}

func (d *UserStorage) GetUser(ctx context.Context, email string) (serviceModels.UserRequest, error) {
	var user databaseModel.User
	err := d.db.GetContext(ctx, &user, `SELECT id, firstname, lastname, password, email, user_type, created_at, updated_at FROM users WHERE email = $1`, email)
	if err != nil {
		return serviceModels.UserRequest{}, err
	}
	return user.ToUserService(), nil
}
