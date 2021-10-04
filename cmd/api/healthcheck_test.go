package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	t.Run("get healthcheck json", func(t *testing.T) {
		app := &application{}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		wantBody := "{\n\t\"status\": \"available\"\n}\n"
		wantStatus := http.StatusOK
		wantContentType := "application/json"

		app.healthcheckHandler(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != wantStatus {
			t.Errorf("status code: got %d want %d", resp.StatusCode, wantStatus)
		}
		if resp.Header.Get("Content-Type") != wantContentType {
			t.Errorf("content-type: got %s\nwant %s", resp.Header.Get("Content-Type"), wantContentType)
		}
		if string(body) != wantBody {
			t.Errorf("body: got\n%s want\n%s", string(body), wantBody)
		}
	})
}
