package common

import (
	"fmt"
	"ginEssential/lxz/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "luckyli",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	fmt.Println("parseToken: ", token)

	return token, claims, err
}
