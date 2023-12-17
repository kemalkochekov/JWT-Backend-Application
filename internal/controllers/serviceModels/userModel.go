package serviceModels

import "time"

// Validate should without empty spaces.
type UserRequest struct {
	ID           int64     `json:"id"`
	Firstname    string    `json:"firstname" validate:"required,min=2,max=100"`
	Lastname     string    `json:"lastname" validate:"required,min=2,max=100"`
	Password     string    `json:"password" validate:"required,min=6,max=100"`
	Email        string    `json:"email" validate:"email,required"`
	Token        string    `json:"token"`
	UserType     string    `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
