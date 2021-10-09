package main

import (
	"context"
	"net/http"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/ilyapetrovMO/WFMTtest/internal/token"
)

type userClaimsKey struct{}

type userClaims struct {
	valid  bool
	userId int
	roleId int
}

func (app *application) managerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userClaimsKey{}).(*userClaims)
		if ok && claims.valid && claims.roleId == db.ROLE_MANAGER {
			next.ServeHTTP(w, r)
		} else {
			app.unauthorizedResponse(w, r, "unauthorized")
		}
	})
}

func (app *application) userOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userClaimsKey{}).(*userClaims)
		if ok && claims.valid && claims.roleId == db.ROLE_USER {
			next.ServeHTTP(w, r)
		} else {
			app.unauthorizedResponse(w, r, "unauthorized")
		}
	})
}

func (app *application) getJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := &userClaims{valid: false}

		tok, err := app.getTokenFromHeader(&r.Header)
		if err == nil {
			c, err := token.ParseJWT(tok)
			if err == nil {
				claims.userId = c.UserId
				claims.roleId = c.RoleId
				claims.valid = true
			}
		}

		ctx := context.WithValue(r.Context(), userClaimsKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
