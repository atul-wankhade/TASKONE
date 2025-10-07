package controller

import (
	"TASKONE/model"
	"TASKONE/repository"
	"TASKONE/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type UserController struct {
	UserRepo repository.UserRepository
	LogRepo  repository.LogRepository
}

func NewUserController(ur repository.UserRepository, lr repository.LogRepository) *UserController {
	return &UserController{
		UserRepo: ur,
		LogRepo:  lr,
	}
}

func (c *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idstr)
	if err != nil || idstr == "" {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"err : ": "Invalid Request"})
		return
	}

	user, err := c.UserRepo.GetByID(id)
	if err != nil {
		utils.JSONResponse(w, http.StatusNotFound, map[string]string{"err : ": "User not found"})
		return
	}

	_ = c.LogRepo.Insert(model.UserLog{UserID: user.ID, Action: "Fethed User", Timestamp: time.Now()}).Error()

	logs, _ := c.LogRepo.GetByUserID(user.ID)

	utils.JSONResponse(w, http.StatusOK, map[string]any{"user": user, "logs": logs})

}

func (c *UserController) GetUserHandlerNew(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idstr)
	if err != nil || idstr == "" {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{"err : ": "Invalid Request"})
		return
	}

	user, err := c.UserRepo.GetByID(id)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return
	}

	if user == nil {
		utils.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	newErr := c.LogRepo.Insert(model.UserLog{UserID: user.ID, Action: "Fetched User", Timestamp: time.Now()})
	if newErr != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{"Err:": "MongoDB query failed"})
		return
	}

	logs, _ := c.LogRepo.GetByUserID(user.ID)

	utils.JSONResponse(w, http.StatusOK, map[string]any{"user": user, "logs": logs})

}
