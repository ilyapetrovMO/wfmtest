package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserModel struct {
	DB *pgxpool.Pool
}

type User struct {
	User_id       int64  `json:"user_id"`
	Username      string `json:"username"`
	Password_hash string `json:"password_hash"`
	Role_id       int64  `json:"role_id"`
}

func (u *UserModel) GetUserByUsrname(ctx context.Context, username string) (*User, error) {
	row := u.DB.QueryRow(ctx, "select user_id, username, password_hash, role_id from users where username=$1", username)
	user := &User{}
	err := row.Scan(&user.User_id, &user.Username, &user.Password_hash, &user.Role_id)

	if err == pgx.ErrNoRows {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserModel) GetUserById(ctx context.Context, userid int) (*User, error) {
	row := u.DB.QueryRow(ctx, "select user_id, username, password_hash, role_id from users where user_id=$1", userid)

	user := &User{}
	err := row.Scan(&user.User_id, &user.Username, &user.Password_hash, &user.Role_id)
	if err == pgx.ErrNoRows {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}
