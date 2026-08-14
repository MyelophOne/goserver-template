package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHealthcheckURL(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")

	if got, want := healthcheckURL(), "http://127.0.0.1:9090/healthz"; got != want {
		t.Fatalf("healthcheckURL() = %q, want %q", got, want)
	}
}

func TestCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/healthz" {
			t.Errorf("path = %q, want /healthz", r.URL.Path)
		}
		if got := r.Header.Get("User-Agent"); got != healthcheckUserAgent {
			t.Errorf("User-Agent = %q, want %q", got, healthcheckUserAgent)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &http.Client{Timeout: time.Second}
	if err := check(client, server.URL+"/healthz"); err != nil {
		t.Fatalf("check() error = %v", err)
	}
}

func TestCheckRejectsUnhealthyStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := &http.Client{Timeout: time.Second}
	err := check(client, server.URL)
	if err == nil {
		t.Fatal("check() error = nil, want unhealthy status error")
	}
	if !strings.Contains(err.Error(), "503 Service Unavailable") {
		t.Fatalf("check() error = %q, want status in error", err)
	}
}
