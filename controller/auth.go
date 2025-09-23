package controller

import (
	"TASKONE/config"
	"TASKONE/utils"
	"encoding/json"
	"net/http"
	"time"

	// "github.com/go-playground/validator"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

var validate = validator.New()

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"Error": "Method not allowed"})
		return
	}

	var loginRequest LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"Error": "Bad Request"})
		return
	}

	err = validate.Struct(loginRequest)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"Error": "validation failed"})
		return
	}

	if loginRequest.Username != "admin" || loginRequest.Password != "admin" {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"Error": "Unauthorized"})
		return
	}
	//create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginRequest.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	secret := []byte(config.AppConfig.JWTSecret)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"Error": "Failed to create token."})
		return
	}

	resp := LoginResponse{
		Token: tokenString,
	}

	utils.JSONResponse(w, http.StatusOK, resp)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]string{"Message": "U have accessed protected endpoint"})

}
