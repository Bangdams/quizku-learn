package model

import "time"

type QuizResponse struct {
	ID         uint      `json:"id" validate:"required"`
	CourseId   uint      `json:"course_id" validate:"required"`
	ClassId    uint      `json:"class_id" validate:"required"`
	QuestionId uint      `json:"question_id" validate:"required"`
	Deadline   uint      `json:"deadline" validate:"required"`
	Status     string    `json:"status" validate:"required"`
	CreatedAt  time.Time `json:"created_at" validate:"required"`
}

type QuizRequest struct {
	CourseId   uint   `json:"course_id" validate:"required"`
	ClassId    uint   `json:"class_id" validate:"required"`
	QuestionId uint   `json:"question_id" validate:"required"`
	Deadline   uint   `json:"deadline" validate:"required"`
	Status     string `json:"status" validate:"required"`
}
