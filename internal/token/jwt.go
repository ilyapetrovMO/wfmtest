package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	dummySecret      = []byte("secret")
	ErrTokenNotValid = errors.New("token could not be validated")
)

type UserClaims struct {
	*jwt.StandardClaims
	User_id int
	Role_id int
}

type User struct {
	User_id int
	Role_id int
}

func NewJWT(user_id, role_id int) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = &UserClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
		user_id,
		role_id,
	}

	tokenString, err := t.SignedString(dummySecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unepxected signing method: %v", token.Header["alg"])
		}

		return dummySecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); !ok {
		userid := claims.User_id
		roleid := claims.Role_id
		return &User{User_id: userid, Role_id: roleid}, nil
	} else {
		return nil, ErrTokenNotValid
	}
}
