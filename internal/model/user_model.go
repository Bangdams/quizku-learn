package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserResponse struct {
	ID        uint   `json:"id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Role      string `json:"role" validate:"required"`
	Image     string `json:"image" validate:"required"`
	CreatedAt string `json:"created_at" validate:"required"`
}

type AdminDashboardReportResponse struct {
	TotalUsers     int64 `json:"total_users" validate:"required"`
	TotalQuizzes   int64 `json:"total_quizzes" validate:"required"`
	TotalClasses   int64 `json:"total_classes" validate:"requied"`
	TotalQuestions int64 `json:"total_questions" validate:"requied"`
}

type LecturerDashboardReportResponse struct {
	TotalUsers   int64 `json:"total_users" validate:"required"`
	TotalQuizzes int64 `json:"total_quizzes" validate:"required"`
	TotalClasses int64 `json:"total_classes" validate:"requied"`
	TotalCourses int64 `json:"total_courses" validate:"requied"`
}

type UserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
	Image    string `json:"image" validate:"required"`
	ClassId  uint   `json:"class_id"`
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
