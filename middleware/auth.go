package middleware

import (
	"TASKONE/config"
	"TASKONE/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"Error": "Missing Authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid Authorization header"})
			return
		}

		tokenstr := parts[1]
		secret := []byte(config.AppConfig.JWTSecret)

		token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKey
			}
			return secret, nil
		})

		if !token.Valid || err != nil {
			utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"Error": "Invalid token"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
