package model

import "github.com/golang-jwt/jwt/v5"

type UserResponse struct {
	ID    uint   `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required"`
	Image string `json:"image" validate:"required"`
}

type UserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
	Image    string `json:"image" validate:"required"`
}

type UpdateUserPasswordRequest struct {
	UpdateUserRequest
}

type UpdateUserRequest struct {
	ID       uint   `json:"id"`
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password"`
	Image    string `json:"image" validate:"required"`
}

type DeleteUserRequest struct {
	ID uint `json:"id" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type TokenPyload struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
