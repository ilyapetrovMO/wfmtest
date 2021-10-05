package main

import (
	"net/http"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/ilyapetrovMO/WFMTtest/internal/token"
)

func (app *application) managerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokstr, err := app.getTokenFromHeader(&r.Header)
		if err != nil {
			app.unauthorizedResponse(w, r, err.Error())
			return
		}

		claims, err := token.ParseJWT(tokstr)
		if err != nil {
			app.unauthorizedResponse(w, r, err.Error())
		}
		if claims.RoleId != db.ROLE_MANAGER {
			app.unauthorizedResponse(w, r, "must be manager")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) userOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokstr, err := app.getTokenFromHeader(&r.Header)
		if err != nil {
			app.unauthorizedResponse(w, r, err.Error())
			return
		}

		claims, err := token.ParseJWT(tokstr)
		switch {
		case err != nil:
			fallthrough
		case claims.RoleId != db.ROLE_USER:
			app.unauthorizedResponse(w, r, "invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
