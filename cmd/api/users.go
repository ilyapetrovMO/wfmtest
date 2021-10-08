package main

import (
	"errors"
	"net/http"

	"github.com/ilyapetrovMO/WFMTtest/internal/db"
	"github.com/ilyapetrovMO/WFMTtest/internal/token"
	"github.com/ilyapetrovMO/WFMTtest/internal/validator"
)

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	input := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := app.readJSON(w, r, input)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	usr, err := app.models.Users.GetByUsername(input.Username)
	if err == db.ErrRecordNotFound {
		app.unauthorizedResponse(w, r, "user not found")
		return
	}

	ok, err := usr.Password.Matches(input.Password)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if !ok {
		app.unauthorizedResponse(w, r, "unauthorized")
	}

	tokenstr, err := token.NewJWT(int(usr.UserId), int(usr.RoleId))
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, &dataJSON{"token": tokenstr})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	input := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := app.readJSON(w, r, input)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	passwd := db.Password{}
	passwd.Set(input.Password)

	user := &db.User{
		Username: input.Username,
		Password: passwd,
		RoleId:   db.ROLE_USER,
	}

	v := validator.New()

	db.ValidateUser(v, user)

	if !v.Valid() {
		app.badRequestResponse(w, r, v.Errors())
		return
	}

	err = app.models.Users.Create(user)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			app.alreadyExistsResponse(w, r, "user exists")
			return
		}
	}

	cart := &db.Cart{UserId: user.UserId}
	err = app.models.Carts.Create(cart)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, &dataJSON{"userId": user.UserId})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}
