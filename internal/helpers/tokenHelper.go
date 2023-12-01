package helpers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type SignedDetails struct {
	Email     string
	Lastname  string
	Firstname string
	UserID    int64
	UserType  string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstname string, lastname string, id int64, userType string) (string, string, error) {
	claims := SignedDetails{
		Email:     email,
		Lastname:  lastname,
		Firstname: firstname,
		UserID:    id,
		UserType:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * 5).Unix(),
		},
	}
	refreshClaims := SignedDetails{
		Email:     email,
		Lastname:  lastname,
		Firstname: firstname,
		UserID:    id,
		UserType:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * 15).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}
func ValidateToken(signedToken string) (*SignedDetails, string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err.Error()
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, fmt.Sprintf("The token is invalid!!!", err.Error())
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Sprintf("The token is expired", err.Error())
	}
	return claims, ""
}
