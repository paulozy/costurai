package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulozy/costurai/pkg"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Name     string `json:"name"`
	Location Location `json:"location"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func NewUser(email, password, name string, location Location) (*User, error) {
	passHash, err := pkg.Encrypt(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(passHash),
		Name:      name,
		Location:  location,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return user, nil
}

func (user *User) UpdateLocation(location Location) {
	user.Location = location
}

func (user *User) UpdateUser(name string, location Location) {
	user.Name = name
	user.Location = location
	user.UpdatedAt = time.Now().Format(time.RFC3339)
}