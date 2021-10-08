package db

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrAlreadyExists  = errors.New("record already exists")
)

const (
	ROLE_MANAGER = 1
	ROLE_USER    = 2
)

type Models struct {
	Users    UserModel
	Products ProductModel
	Orders   OrderModel
	Carts    CartModel
	CartItem CartItemModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		UserModel{db},
		ProductModel{db},
		OrderModel{db},
		CartModel{db},
		CartItemModel{db},
	}
}

func isPgDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return true
		}
	}
	return false
}
