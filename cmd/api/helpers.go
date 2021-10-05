package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type dataJSON map[string]interface{}

var (
	ErrNoAuthorizationHeader = errors.New("authorization header not provided")
	ErrNoToken               = errors.New("no token found in authorization header")
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data *dataJSON) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(dst)
	if err != nil {
		fmt.Printf("ERROR: err in readJSON, %s\n", err)
		return errors.New("could not parse json")
	}

	return nil
}

func (app *application) getTokenFromHeader(h *http.Header) (string, error) {
	tokstr := h.Get("Authorization")
	if tokstr == "" {
		return "", ErrNoAuthorizationHeader
	}

	if !strings.Contains(tokstr, "Bearer") {
		return "", ErrNoToken
	}

	tokstr = strings.TrimPrefix(tokstr, "Bearer")
	tokstr = strings.TrimSpace(tokstr)

	return tokstr, nil
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
