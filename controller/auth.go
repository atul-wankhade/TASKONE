package controller

import (
	"TASKONE/config"
	"TASKONE/model"
	"TASKONE/repository"
	"TASKONE/utils"
	"encoding/json"
	"net/http"
	"time"

	// "github.com/go-playground/validator"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type AuthController struct {
	UserRepo repository.UserRepository
}

func NewAuthController(ur repository.UserRepository) *AuthController {
	return &AuthController{
		UserRepo: ur,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler godoc
// @Summary Login a user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (a *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := a.UserRepo.GetByEmail(loginRequest.Email)
	if err != nil {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"err:": "Invalid Credentials"})
	}

	//compare password with stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)); err != nil {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{"err:": "Invalid Credentials"})
	}

	claims := jwt.MapClaims{
		"sub":   int(user.ID),
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(),
		"email": user.Email,
	}
	//create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte(config.AppConfig.JWTSecret)
	signed, err := token.SignedString(secret)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"Error": "Failed to create token."})
		return
	}

	utils.JSONResponse(w, http.StatusOK, LoginResponse{Token: signed})
}

type registerReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.User true "User registration data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
func (a *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"err": "Method Not Allowed"})
		return
	}

	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Error("invalid request body", zap.Error(err))
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"err": "Invalid Request"})
		return
	}

	if err := validate.Struct(&req); err != nil {
		utils.Logger.Warn("missing required fields", zap.String("email", req.Email))
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"err:": "Validation Failed"})
		return
	}

	h, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Error("failed to hash password", zap.Error(err))
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"err:": "Failed to hash password"})
		return
	}

	u := &model.User{Name: req.Name, Email: req.Email, PasswordHash: string(h)}

	id, err := a.UserRepo.Create(u)
	if err != nil {
		utils.Logger.Error("Failed to create user", zap.Error(err))
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"err:": "Failed to create user"})
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]any{"id": id})

}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]string{"Message": "U have accessed protected endpoint"})

}
