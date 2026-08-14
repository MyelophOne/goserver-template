package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const healthcheckUserAgent = "goserver-healthcheck/1.0"

func healthcheckURL() string {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	return "http://" + net.JoinHostPort("127.0.0.1", port) + "/healthz"
}

func check(client *http.Client, target string) error {
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return fmt.Errorf("create healthcheck request: %w", err)
	}

	req.Header.Set("User-Agent", healthcheckUserAgent)
	req.Close = true

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request health endpoint: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("health endpoint returned %s", response.Status)
	}

	return nil
}

func main() {
	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			Proxy:             nil,
			DisableKeepAlives: true,
		},
	}

	if err := check(client, healthcheckURL()); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "healthcheck failed:", err)
		os.Exit(1)
	}
}
