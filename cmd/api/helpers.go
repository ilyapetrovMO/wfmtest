package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type dataJSON map[string]interface{}

var (
	ErrNoAuthorizationHeader = errors.New("authorization header not provided")
	ErrMalformedHeader       = errors.New("error in authorization header")
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unkown key %s", fieldName)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) getTokenFromHeader(h *http.Header) (string, error) {
	tokstr := h.Get("Authorization")
	if tokstr == "" {
		return "", ErrNoAuthorizationHeader
	}

	if !strings.Contains(tokstr, "Bearer") {
		return "", ErrMalformedHeader
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

func (app *application) readUserIDParam(r *http.Request) (int64, error) {
	userId := r.URL.Query().Get("userId")

	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid userId parameter")
	}

	return id, nil
}

func (app *application) readCartItemIDParam(r *http.Request) (int64, error) {
	cartItemId := r.URL.Query().Get("cartItemId")

	id, err := strconv.ParseInt(cartItemId, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid cartItemId parameter")
	}

	return id, nil
}

func (app *application) readCartIDParam(r *http.Request) (int64, error) {
	cartItemId := r.URL.Query().Get("cartId")

	id, err := strconv.ParseInt(cartItemId, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid cartItemId parameter")
	}

	return id, nil
}
