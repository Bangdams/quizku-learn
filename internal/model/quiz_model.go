package model

import "time"

type QuizResponse struct {
	ID         uint      `json:"id" validate:"required"`
	CourseCode string    `json:"course_code" validate:"required"`
	ClassId    uint      `json:"class_id" validate:"required"`
	QuestionId uint      `json:"question_id" validate:"required"`
	Deadline   time.Time `json:"deadline" validate:"required"`
	Status     string    `json:"status" validate:"required"`
	CreatedAt  time.Time `json:"created_at" validate:"required"`
}

type QuizDashboardResponse struct {
	ID            uint      `json:"id" validate:"required"`
	CourseName    string    `json:"course_name" validate:"required"`
	QuestionCount uint      `json:"question_count" validate:"required"`
	StudentCount  uint      `json:"student_count" validate:"required"`
	CreatedAt     time.Time `json:"created_at" validate:"required"`
}

type QuizRequest struct {
	CourseCode string `json:"course_code" validate:"required"`
	ClassId    uint   `json:"class_id" validate:"required"`
	QuestionId uint   `json:"question_id" validate:"required"`
	Deadline   uint   `json:"deadline" validate:"required"`
	Status     string `json:"status" validate:"required"`
}
