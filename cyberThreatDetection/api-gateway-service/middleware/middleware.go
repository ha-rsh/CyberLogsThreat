package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)


var (
	limiters = make(map[string]*rate.Limiter)
	mu    sync.Mutex
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Auth(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-KEY")
			if key == "" || key != apiKey{
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"success":false,"error":{"code":401,"message":"Unauthorized"}}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RateLimit(requests int, duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)
			limiter := getLimiter(ip, requests, duration)

			if !limiter.Allow() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"success":false,"error":{"code":429,"message":"Rate limit exceeded"}}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getLimiter(key string, requests int, duration time.Duration) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(duration/time.Duration(requests)), requests)
		limiters[key] = limiter
	}

	return limiter
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		println(r.Method, r.URL.Path, duration.String())
	})
}