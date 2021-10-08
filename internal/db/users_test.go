package db

import (
	"context"
	"errors"
	"testing"

	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestUserValidatePasswordPlaintext(t *testing.T) {
	tt := []struct {
		name     string
		Password string
		want     bool
	}{
		{
			"empty Password",
			"",
			false,
		},
		{
			"short Password",
			"1",
			false,
		},
		{
			"long Password",
			"al;sdkfjasdkl;fjasld;kfjasl;dkfjasldkfjasl;dkfjasl;dkfjasldkfjasld;kfjasldkalskdjfasldkjfasld;kfjasldk;fj",
			false,
		},
		{
			"valid Password",
			"12345678ABC",
			true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			v := validator.New()

			ValidatePasswordPlaintext(v, test.Password)
			got := v.Valid()

			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestUserValidateUser(t *testing.T) {
	t.Run("empty hash panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected to panic on empty hash")
			}
		}()

		user := &User{}
		v := validator.New()

		ValidateUser(v, user)
	})

	longUsername := ""
	for i := 0; i <= 500; i++ {
		longUsername += "A"
	}

	testCases := []struct {
		desc string
		user *User
		want int
	}{
		{
			desc: "empty username",
			user: &User{
				Username: "",
				RoleId:   1,
				Password: Password{hash: []byte("1")},
			},
			want: 1,
		}, {
			desc: "long username",
			user: &User{
				Username: longUsername,
				RoleId:   1,
				Password: Password{hash: []byte("1")},
			},
			want: 1,
		}, {
			desc: "0(i.e. empty) role id",
			user: &User{
				Username: "YEP",
				RoleId:   0,
				Password: Password{hash: []byte("1")},
			},
			want: 1,
		}, {
			desc: "bad role id (must be 1 or 2)",
			user: &User{
				Username: "YEP",
				RoleId:   3,
				Password: Password{hash: []byte("1")},
			},
			want: 1,
		}, {
			desc: "all fields erronious",
			user: &User{
				Username: "",
				RoleId:   0,
				Password: Password{hash: []byte("1")},
			},
			want: 2,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			v := validator.New()

			ValidateUser(v, tC.user)

			if len(v.Errors()) != tC.want {
				t.Errorf("got errors: %v want len %d", v.Errors(), tC.want)
			}
		})
	}
}

func TestUserGetByUsername(t *testing.T) {
	// Expect to have user1 and manager1 seeded in test DB
	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	unexpectedErr(t, err)
	defer dbpool.Close()

	um := &UserModel{dbpool}

	testCases := []struct {
		desc     string
		username string
		want     string
		err      error
	}{
		{
			desc:     "get user1",
			username: "user1",
			want:     "user1",
		}, {
			desc:     "get manager1",
			username: "manager1",
			want:     "manager1",
		}, {
			desc:     "try getting an nonexisting user",
			username: "!-234_0",
			err:      ErrRecordNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := um.GetByUsername(tC.username)
			if err != nil {
				switch {
				case tC.err != nil && errors.Is(err, tC.err):
					return
				default:
					t.Fatalf("unexpected error %s", err)
					return
				}
			}

			if got.Username != tC.want {
				t.Errorf("got %s want %s", got.Username, tC.want)
			}
		})
	}
}

func TestUserGetById(t *testing.T) {
	// Expect to have user1 and manager1 seeded in test DB
	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	unexpectedErr(t, err)
	defer dbpool.Close()

	um := &UserModel{dbpool}

	testCases := []struct {
		desc   string
		userId int
		want   string
		err    error
	}{
		{
			desc:   "get user1",
			userId: 1,
			want:   "user1",
		}, {
			desc:   "get manager1",
			userId: 2,
			want:   "manager1",
		}, {
			desc:   "try getting an nonexisting user",
			userId: -1,
			err:    ErrRecordNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := um.GetById(tC.userId)
			if err != nil {
				switch {
				case tC.err != nil && errors.Is(err, tC.err):
					return
				default:
					t.Fatalf("unexpected error %s", err)
					return
				}
			}

			if got.Username != tC.want {
				t.Errorf("got %s want %s", got.Username, tC.want)
			}
		})
	}
}

func TestUserCreate(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	unexpectedErr(t, err)
	defer dbpool.Close()
	defer dbpool.Exec(context.Background(), `DELETE FROM users WHERE username LIKE '%TEST'`)

	t.Run("create user", func(t *testing.T) {
		um := &UserModel{dbpool}

		us := &User{
			Username: "testCreateTEST",
			RoleId:   1,
			Password: Password{hash: []byte("1")},
		}

		err := um.Create(us)
		unexpectedErr(t, err)
	})

	t.Run("create user twice", func(t *testing.T) {
		um := &UserModel{dbpool}

		us := &User{
			Username: "createTwiceTEST",
			RoleId:   1,
			Password: Password{hash: []byte("1")},
		}

		err := um.Create(us)
		unexpectedErr(t, err)

		err = um.Create(us)
		if err == nil {
			t.Fatal("expected error")
		}

		if !errors.Is(err, ErrAlreadyExists) {
			unexpectedErr(t, err)
		}
	})
}

func unexpectedErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
}
