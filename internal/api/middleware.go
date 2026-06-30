package api

import (
	"net/http"
	"os"
)

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwt string

			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}
			var valid bool
			valid = ValidateToken(jwt)

			if !valid {
				writeJSON(w, map[string]interface{}{"error": "Authentification required"}, http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
