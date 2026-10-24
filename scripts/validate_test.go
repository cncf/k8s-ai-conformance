//go:build validate

package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func withFastBackoff(t *testing.T) {
	t.Helper()
	oldBase, oldMax := baseBackoff, maxBackoff
	baseBackoff, maxBackoff = time.Millisecond, 5*time.Millisecond
	t.Cleanup(func() { baseBackoff, maxBackoff = oldBase, oldMax })
}

func TestValidateURLRetries429ThenSucceeds(t *testing.T) {
	withFastBackoff(t)
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) <= 2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if err := validateURL(srv.URL); err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Fatalf("expected 3 requests, got %d", got)
	}
}

func TestValidateURLPersistent429IsRateLimited(t *testing.T) {
	withFastBackoff(t)
	var methods = map[string]int{}
	var mu = make(chan struct{}, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu <- struct{}{}
		methods[r.Method]++
		<-mu
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	err := validateURL(srv.URL)
	if !isRateLimited(err) {
		t.Fatalf("expected rate-limited error, got %v", err)
	}
	if methods[http.MethodHead] != maxRetries+1 || methods[http.MethodGet] != maxRetries+1 {
		t.Fatalf("expected %d HEAD and GET attempts, got %v", maxRetries+1, methods)
	}
}

func TestValidateURLNotFoundIsError(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	err := validateURL(srv.URL)
	if err == nil || isRateLimited(err) {
		t.Fatalf("expected a hard error for 404, got %v", err)
	}
}
