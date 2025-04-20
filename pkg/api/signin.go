package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// SignInReq — схема входящего JSON
type SignInReq struct {
	Password string `json:"password"`
}

// SignInResp — схема ответа
type SignInResp struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// POST /api/signin
func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, SignInResp{Error: "invalid JSON"})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		// если пароль не задан — сразу впускаем
		writeJSON(w, SignInResp{Token: ""})
		return
	}
	if req.Password != pass {
		writeJSON(w, SignInResp{Error: "invalid password"})
		return
	}

	// создаём JWT с минимальным payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
	})
	tokenStr, err := token.SignedString([]byte(pass))
	if err != nil {
		writeJSON(w, SignInResp{Error: "could not create token"})
		return
	}

	// ставим куку token
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Path:     "/",
		Expires:  time.Now().Add(8 * time.Hour),
		HttpOnly: true,
	})

	writeJSON(w, SignInResp{Token: tokenStr})
}
