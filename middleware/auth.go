package middleware

import (
	"TASKONE/config"
	"TASKONE/utils"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type ctxKey string

var UserIDKey ctxKey = "userID"

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
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrInvalidKey
			}
			return secret, nil
		})

		if !token.Valid || err != nil {
			utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"Error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//expaect sub claims as integer , we stored it as int
			if sub, ok := claims["sub"].(float64); ok {
				uid := int(sub)
				ctx := context.WithValue(r.Context(), UserIDKey, uid)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"err:": "Invalid token claims"})
	})
}
