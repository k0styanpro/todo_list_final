package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// auth — мидлварь, оборачивает API‑эндпоинты.
// Если TODO_PASSWORD не задан, пропускает без проверки.
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}
			tokenStr := cookie.Value

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				// проверяем метод подписи
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(pass), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}
