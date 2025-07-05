// middleware/ratelimit.go
// Contains the IP-based rate limiting middleware.
package middleware

import (
	"net/http"
	"sync"
	"time"

	"your_username/budget-tracker/config"

	"golang.org/x/time/rate"
)

// client represents a user of the API, identified by IP address.
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients map[string]*client
	mu      sync.Mutex
	cfg     config.RateLimiterConfig
)

// InitRateLimiter initializes the rate limiter with configuration
// and starts the background cleanup goroutine.
func InitRateLimiter(limiterConfig config.RateLimiterConfig) {
	cfg = limiterConfig
	clients = make(map[string]*client)

	// Start a background goroutine to remove old entries from the clients map.
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

// getClient returns the rate limiter for the provided IP address.
func getClient(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	c, exists := clients[ip]
	if !exists {
		// Create a new rate limiter and add it to the clients map.
		limiter := rate.NewLimiter(rate.Limit(cfg.RPS), cfg.Burst)
		clients[ip] = &client{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	// Update the last seen time for the client.
	c.lastSeen = time.Now()
	return c.limiter
}

// RateLimiter is the middleware function.
func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If rate limiting is disabled, just call the next handler.
		if !cfg.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		// Get the IP address for the current user.
		ip := r.RemoteAddr
		limiter := getClient(ip)

		// Check if the request is allowed.
		if !limiter.Allow() {
			http.Error(w, "API rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
