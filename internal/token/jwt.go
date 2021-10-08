package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	dummySecret      = []byte("secret")
	ErrTokenNotValid = errors.New("token could not be validated")
)

type UserClaims struct {
	*jwt.StandardClaims
	UserId int
	RoleId int
}

func NewJWT(userId, roleId int) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = &UserClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
		userId,
		roleId,
	}

	tokenString, err := t.SignedString(dummySecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return dummySecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	} else {
		return nil, ErrTokenNotValid
	}
}
