package db

import (
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

const (
	ROLE_MANAGER = 1
	ROLE_USER    = 2
)

type Models struct {
	Users    UserModel
	Products ProductModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		UserModel{db},
		ProductModel{db},
	}
}
