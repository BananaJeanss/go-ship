package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func FormatHandler(w http.ResponseWriter, r *http.Request) {
	// just sends request to playground with compliance

	// 1 format per 1 second
	ip := getIP(r)
	if !checkRateLimit(ip, clientLastFormat, 1*time.Second) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"Error": "Rate limit exceeded. Please wait 1 second before formatting again."}`))
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

	// The format endpoint expects urlencoded form data
	formData := url.Values{}
	formData.Set("imports", "true")
	formData.Set("body", string(bodyBytes))

	req, err := http.NewRequest("POST", "https://go.dev/_/fmt?backend=", strings.NewReader(formData.Encode()))
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
		http.Error(w, "Failed to format code", http.StatusInternalServerError)
		fmt.Printf("debug data response: %s\n", resp.Status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
