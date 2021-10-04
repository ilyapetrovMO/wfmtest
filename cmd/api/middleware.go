package main

import (
	"net/http"
	"strings"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/ilyapetrovMO/WFMTtest/internal/token"
)

func (app *application) managerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokstr := r.Header.Get("Authorization")

		if !strings.Contains(tokstr, "Bearer") {
			app.unauthorizedResponse(w, r)
			return
		}

		tokstr = strings.TrimPrefix(tokstr, "Bearer")
		tokstr = strings.TrimSpace(tokstr)

		usrclaims, err := token.ParseJWT(tokstr)
		switch {
		case err != nil:
			fallthrough
		case usrclaims.Role_id != db.ROLE_MANAGER:
			app.unauthorizedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
