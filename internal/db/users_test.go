package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestGetUser(t *testing.T) {
	var connstring = "postgres://postgres:postgres@127.0.0.1:5432/wfmtest"

	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	unexpectedErr(t, err)
	defer dbpool.Close()

	want := &User{
		User_id:       1,
		Username:      "user1",
		Role_id:       ROLE_USER,
		Password_hash: "$2a$14$ymJHFkT1IO2PxAovxD83j.WNGpf5SqCP2zV9x/UoVzCMO6mvxDr4W",
	}

	um := &UserModel{dbpool}
	got, err := um.GetUser(context.Background(), "user1")
	unexpectedErr(t, err)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

func unexpectedErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
}
