package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	helloHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	buf := new(strings.Builder)
	_, err := io.Copy(buf, resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Check that the response body is exactly what we expect
	body := buf.String()
	expected := "Hello from moshi-moshi"
	if body != expected {
		t.Errorf("expected body %q, got %q", expected, body)
	}
}

func TestParamsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/params?foo=bar", nil)
	req.Header.Set("User-Agent", "GoTest")
	req.Header.Set("X-Test-Header", "test-value")
	// Use a proper host:port format for RemoteAddr
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	paramsHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var payload map[string]interface{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&payload); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	client, ok := payload["client"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing or invalid 'client' key in response")
	}
	server, ok := payload["server"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing or invalid 'server' key in response")
	}

	// Check that we get the IP without the port
	if client["source_ip"] != "127.0.0.1" {
		t.Errorf("expected source_ip '127.0.0.1', got %v", client["source_ip"])
	}
	if client["user_agent"] != "GoTest" {
		t.Errorf("expected user_agent 'GoTest', got %v", client["user_agent"])
	}
	if client["http_verb"] != "GET" {
		t.Errorf("expected http_verb 'GET', got %v", client["http_verb"])
	}
	if client["requested_path"] != "/params" {
		t.Errorf("expected requested_path '/params', got %v", client["requested_path"])
	}
	if client["query_string"] != "foo=bar" {
		t.Errorf("expected query_string 'foo=bar', got %v", client["query_string"])
	}

	// Use type assertion with ok check for headers
	headers, ok := client["headers"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing or invalid 'headers' key in client")
	}

	// Check for test header - it might be case-insensitive in the response
	found := false
	testHeaderVal, ok := headers["X-Test-Header"]
	if ok && testHeaderVal == "test-value" {
		found = true
	}
	// Try lowercase version too
	testHeaderVal, ok = headers["x-test-header"]
	if ok && testHeaderVal == "test-value" {
		found = true
	}

	if !found {
		t.Errorf("expected X-Test-Header 'test-value', but header not found or has wrong value")
	}

	if server["listen_port"] != float64(8080) { // JSON numbers are float64
		t.Errorf("expected listen_port 8080, got %v", server["listen_port"])
	}
}
