package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type dataJSON map[string]interface{}

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
		return errors.New("could not parse json")
	}

	return nil
}
