package token

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

var dummySecret = []byte("secret")

type UserClaims struct {
	Username string
	Role_id  int
}

func NewJWT(username string, role_id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role_id":  role_id,
	})

	tokenString, err := token.SignedString(dummySecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (*UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unepxected signing method: %v", token.Header["alg"])
		}

		return dummySecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		roleid := claims["role_id"].(float64)
		return &UserClaims{Username: username, Role_id: int(roleid)}, nil
	} else {
		return nil, err
	}
}
