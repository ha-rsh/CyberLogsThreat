package middleware

import (
	"cybersecuritySystem/shared/auth"
	"cybersecuritySystem/shared/logger"
	"net/http"
	"strings"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/auth/login") ||
			strings.HasSuffix(r.URL.Path, "/auth/register") {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"error":{"code":401,"message":"Missing authorization token"}}`))
			return
		}

		var tokenString string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"error":{"code":401,"message":"Invalid authorization format. Use: Bearer <token>"}}`))
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			logger.Warn("Invalid token: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"success":false,"error":{"code":401,"message":"Invalid or expired token"}}`))
			return
		}

		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Username", claims.Username)
		r.Header.Set("X-User-Role", claims.Role)

		logger.Debug("Authenticated: %s (role: %s)", claims.Username, claims.Role)

		next.ServeHTTP(w, r)
	})
}