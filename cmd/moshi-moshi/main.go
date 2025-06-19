package main

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Set up slog logger
	logFormat := os.Getenv("MOSHI_LOG_FORMAT")
	var handler slog.Handler
	switch logFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, nil)
	default:
		handler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/params", paramsHandler)
	http.HandleFunc("/health", healthcheckHandler)
	http.HandleFunc("/healthcheck", healthcheckHandler)

	slog.Info("Starting server", "addr", ":8080", "log_format", logFormat)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Server failed", "error", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request received", "path", r.URL.Path, "method", r.Method, "remote", r.RemoteAddr)
	w.Write([]byte("Hello from moshi-moshi"))
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Request received", "path", r.URL.Path, "method", r.Method, "remote", r.RemoteAddr)
	w.Write([]byte("WORKING"))
}

func paramsHandler(w http.ResponseWriter, r *http.Request) {
	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}

	// Collect headers
	headers := make(map[string]string)
	for k, v := range r.Header {
		headers[k] = strings.Join(v, ", ")
	}

	// Server info
	hostname, _ := os.Hostname()
	serverInfo := map[string]interface{}{
		"hostname":    hostname,
		"listen_ip":   "0.0.0.0",
		"listen_port": 8080,
		"protocol":    r.Proto,
	}

	clientInfo := map[string]interface{}{
		"source_ip":      clientIP,
		"user_agent":     r.UserAgent(),
		"http_verb":      r.Method,
		"requested_path": r.URL.Path,
		"query_string":   r.URL.RawQuery,
		"headers":        headers,
	}

	payload := map[string]interface{}{
		"client": clientInfo,
		"server": serverInfo,
	}

	slog.Info("Params endpoint hit", "client", clientInfo, "server", serverInfo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
