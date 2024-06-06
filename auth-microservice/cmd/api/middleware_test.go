package main

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nurtikaga/internal/data"
)

func TestRecoverPanic(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.recoverPanic(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})))
	defer ts.Close()

	code, _, _ := ts.get(t, "/")

	if code != http.StatusInternalServerError {
		t.Errorf("want %d; got %d", http.StatusInternalServerError, code)
	}
}

func TestRateLimit(t *testing.T) {
	app := newTestApplication(t)
	app.config.limiter.rps = 2
	app.config.limiter.burst = 4
	app.config.limiter.enabled = true

	ts := newTestServer(t, app.rateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})))
	defer ts.Close()

	for i := 0; i < 10; i++ {
		code, _, body := ts.get(t, "/")
		if i < 6 && code != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, code)
		}
		if i >= 6 && code != http.StatusTooManyRequests {
			t.Errorf("want %d; got %d", http.StatusTooManyRequests, code)
		}
		if i < 6 && strings.TrimSpace(body) != "OK" {
			t.Errorf("want body to equal %q", "OK")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestRequireAuthenticatedUser(t *testing.T) {
	app := newTestApplication(t)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	ts := newTestServer(t, app.requireAuthenticatedUser(next))
	defer ts.Close()

	t.Run("authenticated", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/", nil)
		ctx := context.WithValue(req.Context(), userContextKey, &data.UserInfo{ID: 1})
		req = req.WithContext(ctx)

		rs, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer rs.Body.Close()

		if rs.StatusCode != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
		}
	})

	t.Run("unauthenticated", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/", nil)

		rs, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer rs.Body.Close()

		if rs.StatusCode != http.StatusUnauthorized {
			t.Errorf("want %d; got %d", http.StatusUnauthorized, rs.StatusCode)
		}
	})
}

func TestRequireActivatedUser(t *testing.T) {
	app := newTestApplication(t)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	ts := newTestServer(t, app.requireActivatedUser(next))
	defer ts.Close()

	t.Run("activated", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/", nil)
		ctx := context.WithValue(req.Context(), userContextKey, &data.UserInfo{ID: 1, Activated: true})
		req = req.WithContext(ctx)

		rs, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer rs.Body.Close()

		if rs.StatusCode != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
		}
	})

	t.Run("not activated", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/", nil)
		ctx := context.WithValue(req.Context(), userContextKey, &data.UserInfo{ID: 1, Activated: false})
		req = req.WithContext(ctx)

		rs, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer rs.Body.Close()

		if rs.StatusCode != http.StatusForbidden {
			t.Errorf("want %d; got %d", http.StatusForbidden, rs.StatusCode)
		}
	})
}
