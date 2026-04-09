package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// https://www.alexedwards.net/blog/how-to-rate-limit-http-requests
var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Hour/5), 5)
		visitors[ip] = limiter
	}

	return limiter
}

func addVisitor(ip string) {
	mu.Lock()
	defer mu.Unlock()
	visitors[ip] = rate.NewLimiter(rate.Every(time.Hour/5), 5)
}

// shoutout ai.hackclub.com
func GenerateMeAIdeaHandler(w http.ResponseWriter, r *http.Request) {
	AI_BASE_URL := "https://ai.hackclub.com/proxy/v1/chat/completions"
	AI_API_TOKEN := os.Getenv("AI_API_TOKEN")
	AI_MODEL := "qwen/qwen3-32b"

	if AI_BASE_URL == "" || AI_API_TOKEN == "" {
		http.Error(w, "AI API not configured", 500)
		return
	}

	prompt := "hai"

	// first off, ratelimit of 5 reqs per hour per ip
	// if there's a header e.g. Xforwarded or cloudflare or whatever use that instead
	HeadersToCheck := []string{"X-Forwarded-For", "CF-Connecting-IP", "X-Real-IP"}
	var ip string
	for _, header := range HeadersToCheck {
		if r.Header.Get(header) != "" {
			ip = r.Header.Get(header)
			break
		}
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	limiter := getVisitor(ip)
	if !limiter.Allow() {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// cook the request to ai.hackclub.com
	req, err := http.NewRequest("POST", AI_BASE_URL, nil)
	if err != nil {
		http.Error(w, "Failed to create AI request", 500)
		return
	}
	req.Header.Set("Authorization", "Bearer "+AI_API_TOKEN)
	req.Header.Set("Content-Type", "application/json")
	req.Body =
		func() io.ReadCloser {
			b, _ := json.Marshal(map[string]any{
				"model": AI_MODEL,
				"messages": []map[string]string{
					{"role": "user", "content": prompt},
				},
			})
			return io.NopCloser(bytes.NewReader(b))
		}()

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send AI request", 500)
		return
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read AI response", 500)
		return
	}

	// get the answer from the response
	var aiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &aiResp); err != nil {
		http.Error(w, "Failed to parse AI response", 500)
		return
	}

	if len(aiResp.Choices) == 0 || aiResp.Choices[0].Message.Content == "" {
		http.Error(w, "AI response missing content", 500)
		return
	}

	// everything passed, add to ratelimit and send response
	addVisitor(ip)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(aiResp.Choices[0].Message.Content))
}
