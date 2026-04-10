package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RunHandler(w http.ResponseWriter, r *http.Request) {
	// just sends run request to playground with compliance

	// 1 run per 2 seconds
	ip := getIP(r)
	if !checkRateLimit(ip, clientLastRun, 2*time.Second) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"Errors": "Rate limit exceeded. Please wait 2 seconds before running again."}`))
		return
	}

	// only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// only expect raw text
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusUnsupportedMediaType)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// The compile endpoint expects urlencoded form data
	formData := url.Values{}
	formData.Set("version", "2")
	formData.Set("body", string(bodyBytes))
	formData.Set("withVet", "true")

	req, err := http.NewRequest("POST", "https://go.dev/_/compile?backend=", strings.NewReader(formData.Encode()))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Go Ship! (https://goship.dino.icu)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to run code", http.StatusInternalServerError)
		fmt.Printf("debug data response: %s\n", resp.Status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
