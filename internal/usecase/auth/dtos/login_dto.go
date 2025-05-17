package dtos

import "github.com/paulozy/costurai/internal/entity"

type AuthenticationInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthDressmakerOutput struct {
	Token string            `json:"token"`
	User  entity.Dressmaker `json:"user"`
}

type AuthUserOutput struct {
	Token string      `json:"token"`
	User  entity.User `json:"user"`
}
