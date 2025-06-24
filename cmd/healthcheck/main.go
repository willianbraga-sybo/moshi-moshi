package main

import (
    "net/http"
    "net/url"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        os.Exit(1)
    }
    rawURL := os.Args[1]
    parsedURL, err := url.ParseRequestURI(rawURL)
    if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
        os.Exit(1)
    }
    resp, err := http.Get(rawURL)
    if err != nil || resp.StatusCode >= 400 {
        os.Exit(1)
    }
}
