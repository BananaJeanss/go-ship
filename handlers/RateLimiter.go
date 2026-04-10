package handlers

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	clientLastRun    = make(map[string]time.Time)
	clientLastFormat = make(map[string]time.Time)
	rlMutex          sync.Mutex
)

// getIP safely extracts the actual IP address
func getIP(r *http.Request) string {
	if header := r.Header.Get("X-Forwarded-For"); header != "" {
		ips := strings.Split(header, ",")
		return strings.TrimSpace(ips[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// checkRateLimit returns true if the request is allowed, false if it should be rate limited
func checkRateLimit(ip string, limitsMap map[string]time.Time, delay time.Duration) bool {
	rlMutex.Lock()
	defer rlMutex.Unlock()

	lastReq, exists := limitsMap[ip]
	if !exists || time.Since(lastReq) > delay {
		// isn't ratelimited
		limitsMap[ip] = time.Now() // update last request time
		return true
	}
	// is ratelimited
	return false
}
