package db

import (
	"context"
	"errors"
	"time"

	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *pgxpool.Pool
}

type Password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	UserId    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Password  Password  `json:"-"`
	RoleId    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	Version   int64     `json:"-"`
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(password != "", "password", "must be non-empty string")
	v.Check(len(password) <= 72, "password", "must not be longer than 72 bytes")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 15, "username", "must not be longer than 15 bytes")

	v.Check(user.RoleId != 0, "role_id", "must be provided")
	v.Check(user.RoleId == ROLE_MANAGER || user.RoleId == ROLE_USER, "role_id", "invalid role id")

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("no password hash")
	}
}

func (u *UserModel) Create(user *User) error {
	sql := `
		INSERT INTO users (username, password_hash, role_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id, username, role_id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	args := []interface{}{
		&user.Username,
		&user.Password.hash,
		&user.RoleId,
		time.Now(),
	}

	err := u.DB.QueryRow(ctx, sql, args...).Scan(
		&user.UserId,
		&user.Username,
		&user.RoleId,
		&user.CreatedAt,
		&user.Version,
	)
	if err != nil {
		if isPgDuplicateError(err) {
			return ErrAlreadyExists
		}
		return err
	}

	return nil
}

func (u *UserModel) GetByUsername(username string) (*User, error) {
	sql := `
	SELECT user_id, username, password_hash, role_id, created_at, version
	FROM users
	WHERE username = $1 AND deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	user := &User{}
	err := u.DB.QueryRow(ctx, sql, username).Scan(
		&user.UserId,
		&user.Username,
		&user.Password.hash,
		&user.RoleId,
		&user.CreatedAt,
		&user.Version,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserModel) GetById(userid int) (*User, error) {
	sql := `
	SELECT user_id, username, password_hash, role_id, created_at, version
	FROM users
	WHERE user_id=$1 AND deleted_at IS NULL`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	user := &User{}

	err := u.DB.QueryRow(ctx, sql, userid).Scan(
		&user.UserId,
		&user.Username,
		&user.Password.hash,
		&user.RoleId,
		&user.CreatedAt,
		&user.Version,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRecordNotFound
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserModel) Delete(user *User) error {
	sql := `
	UPDATE users
	SET deleted_at=$1, version=version+1
	WHERE user_id=$2 AND version=$3
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := u.DB.QueryRow(ctx, sql, time.Now(), user.UserId, user.Version).Scan(&user.Version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRecordNotFound
		}

		return err
	}

	return nil
}
