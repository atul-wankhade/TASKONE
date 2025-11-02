package controller

import (
	"TASKONE/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	createFn     func(user *model.User) (int64, error)
	GetByEmailFn func(email string) (*model.User, error)
}

func (m *mockUserRepo) Create(user *model.User) (int64, error) {
	return m.createFn(user)
}

func (m *mockUserRepo) GetByEmail(email string) (*model.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(email)
	}
	return nil, errors.New("Not implemented")
}

func (m *mockUserRepo) GetByID(id int) (*model.User, error) {
	return nil, nil
}

func TestLogin(t *testing.T) {
	mockRepo := &mockUserRepo{

		GetByEmailFn: func(email string) (*model.User, error) {
			hashedPasswd, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
			return &model.User{
				ID:           1,
				Name:         "testuser",
				Email:        email,
				PasswordHash: string(hashedPasswd),
			}, nil
		},
	}

	controller := &AuthController{UserRepo: mockRepo}

	type Config struct {
		JWTSecret string
	}
	AppConfig := &Config{
		JWTSecret: "testsecret",
	}
	fmt.Print(AppConfig)

	data := map[string]string{
		"email":    "test@user.com",
		"password": "test123",
	}

	body, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	controller.LoginHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 ok but got %d", rec.Code)
	}

}
