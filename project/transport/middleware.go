package transport

import (
	"context"
	"net/http"
	"strings"
	"students-api/project/authentication"
	"time"

	log "github.com/sirupsen/logrus"
)

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware - a handy middleware function that logs out incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
			}).
			Info("handled request")
		next.ServeHTTP(w, r)
	})
}

func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("Unauthorized Request.")
			return
		}

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("Error parsing Authorization Header.")
			return
		}
		claims, err := authentication.ValidateJWT(authHeaderParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Error("Failed to validate token.")
			return
		} else {
			log.Info(claims)
			customClaims := claims.UserID
			log.Info(customClaims)
			ctx := context.WithValue(r.Context(), authentication.UserIDKey, customClaims)
			r = r.WithContext(ctx)
			original(w, r)
		}

	}
}
