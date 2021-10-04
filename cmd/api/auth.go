package main

import (
	"fmt"
	"net/http"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/ilyapetrovMO/WFMTtest/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string
	Password string
}

func (app *application) authenticationHandler(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := app.readJSON(r, creds)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "malformed json")
	}

	usr, err := app.models.Users.GetUser(r.Context(), creds.Username)
	if err == db.ErrRecordNotFound {
		app.unauthorizedResponse(w, r)
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password_hash), []byte(creds.Password))
	if err != nil {
		app.unauthorizedResponse(w, r)
	}

	tokenstr, err := token.NewJWT(creds.Username, int(usr.Role_id))
	if err != nil {
		fmt.Printf("error creating token: %s", err)
		app.internalServerErrorResponse(w, r)
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"token": tokenstr})
	if err != nil {
		fmt.Printf("error creating token: %s", err)
		app.internalServerErrorResponse(w, r)
	}
}
