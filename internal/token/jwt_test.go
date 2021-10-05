package token

import "testing"

func TestJWT(t *testing.T) {
	t.Run("create JWT", func(t *testing.T) {
		wantUserId := 2
		wantRoleId := 1

		tokenString, err := NewJWT(wantUserId, wantRoleId)
		unexpectedErr(t, err)

		got, err := ParseJWT(tokenString)
		unexpectedErr(t, err)

		if got.UserId != wantUserId {
			t.Errorf("username: got %d want %d", got.UserId, wantUserId)
		}
		if got.RoleId != wantRoleId {
			t.Errorf("role_id: got %d want %d", got.RoleId, wantRoleId)
		}
	})
}

func unexpectedErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
}
