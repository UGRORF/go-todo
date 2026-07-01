package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	HashedPassword string `json:"password"`
	jwt.RegisteredClaims
}

func SinginHandler(w http.ResponseWriter, r *http.Request) {
	storedPassword := os.Getenv("TODO_PASSWORD")
	if storedPassword == "" {
		writeError(w, "authentication not configured", http.StatusInternalServerError)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Password != storedPassword {
		writeError(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	token, err := CreateToken(storedPassword)
	if err != nil {
		writeError(w, "failed to create Token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{"token": token}, http.StatusOK)
}

func CreateToken(password string) (string, error) {
	hashedPassword := HashPassword(password)
	claims := Claims{
		HashedPassword: hashedPassword,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(getJWTSecret()))
}

func ValidateToken(tokenString string) bool {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return true
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(getJWTSecret()), nil
	})
	if err != nil || !token.Valid {
		return false
	}

	return claims.HashedPassword == HashPassword(password)
}

func HashPassword(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func getJWTSecret() string {
	secret := os.Getenv("TODO_JWT_SECRET")
	if secret == "" {
		secret = "default-secret"
	}
	return secret
}
