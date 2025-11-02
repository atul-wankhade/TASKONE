package controller

import (
	"TASKONE/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterUSer(t *testing.T) {
	mockRepo := &mockUserRepo{
		createFn: func(user *model.User) (int64, error) {
			if user.Email == "exists@example.com" {
				return int64(1), errors.New("User already exists")
			}
			return int64(0), nil
		},
	}

	controller := &AuthController{UserRepo: mockRepo}

	body, _ := json.Marshal(map[string]string{
		"name":     "test User",
		"email":    "test@user.com",
		"password": "test123",
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))

	rec := httptest.NewRecorder()

	controller.RegisterHandler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 but got /%d", rec.Code)
	}
}
