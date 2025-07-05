package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	pass := os.Getenv("TODO_PASSWORD")
	hash := sha256.Sum256([]byte(pass))
	return hash[:]
}

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Auth required", http.StatusUnauthorized)
				return
			}
			tokenStr := cookie.Value

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return getJWTSecret(), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Auth required", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Auth required", http.StatusUnauthorized)
				return
			}

			exp, ok := claims["exp"].(float64)
			if !ok || int64(exp) < time.Now().UTC().Unix() {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}

			hash, ok := claims["hash"].(string)
			sum := sha256.Sum256([]byte(pass))
			expectedHash := hex.EncodeToString(sum[:])
			if !ok || hash != expectedHash {
				http.Error(w, "Auth required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}
