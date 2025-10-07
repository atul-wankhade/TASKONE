package repository

import (
	"TASKONE/model"
	"database/sql"
)

type UserRepository interface {
	GetByID(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(u *model.User) (int64, error)
}

type userRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{DB: db}
}

func (r *userRepo) GetByID(id int) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRow("SELECT id, name, email, password_hash FROM users WHERE id=?;", id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRow("SELECT id, name, email, password_hash FROM users WHERE email=?", email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(u *model.User) (int64, error) {
	res, err := r.DB.Exec("INSERT INTO users (name,email,password_hash) VALUES(?,?,?)", u.Name, u.Email, u.PasswordHash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
