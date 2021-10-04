package token

import "testing"

func TestJWT(t *testing.T) {
	t.Run("create JWT", func(t *testing.T) {
		wantUser := "user1"
		wantRoleId := 1

		tokenString, err := NewJWT("user1", 1)
		unexpectedErr(t, err)

		got, err := ParseJWT(tokenString)
		unexpectedErr(t, err)

		if got.Username != wantUser {
			t.Errorf("username: got %s want %s", got.Username, wantUser)
		}
		if got.Role_id != wantRoleId {
			t.Errorf("role_id: got %d want %d", got.Role_id, wantRoleId)
		}
	})
}

func unexpectedErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
}
