package main

import (
	"net/http"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, _ := ts.get(t, "/v1/healthCheck")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
}
