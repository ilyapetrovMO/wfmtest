package main

import (
	"fmt"
	"net/http"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/ilyapetrovMO/WFMTtest/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *application) authenticationHandler(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := app.readJSON(r, creds)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		app.errorResponse(w, r, http.StatusBadRequest, "could not read JSON body")
		return
	}

	usr, err := app.models.Users.GetUserByUsrname(r.Context(), creds.Username)
	if err == db.ErrRecordNotFound {
		app.unauthorizedResponse(w, r, "user not found")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(creds.Password))
	if err != nil {
		app.unauthorizedResponse(w, r, "incorrect password")
		return
	}

	tokenstr, err := token.NewJWT(int(usr.UserId), int(usr.RoleId))
	if err != nil {
		fmt.Printf("error creating token: %s", err)
		app.internalServerErrorResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"token": tokenstr})
	if err != nil {
		fmt.Printf("error creating token: %s", err)
		app.internalServerErrorResponse(w, r)
		return
	}
}
