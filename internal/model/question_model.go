package model

import "time"

type QuestionResponse struct {
	ID            uint      `json:"id" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	QuestionCount uint      `json:"question_count" validate:"required"`
	Duration      uint      `json:"duration" validate:"required"`
	CourseCode    string    `json:"course_code" validate:"required"`
	UserId        uint      `json:"user_id" validate:"required"`
	CreatedAt     time.Time `json:"created_at" validate:"required"`
}

type QuestionWithCourseUserResponse struct {
	ID           uint      `json:"id" validate:"required"`
	QuestionName string    `json:"question_name" validate:"required"`
	CourseName   string    `json:"course_name" validate:"required"`
	LecturerName string    `json:"lecturer_name" validate:"required"`
	CreatedAt    time.Time `json:"created_at" validate:"required"`
}

type QuestionRequest struct {
	Name           string `json:"name" validate:"required"`
	QuestionCount  uint   `json:"question_count" validate:"required"`
	Duration       uint   `json:"duration" validate:"required"`
	CourseCode     string `json:"course_code" validate:"required"`
	UserId         uint
	QuestionDetail []QuestionDetailRequest `json:"question_detail" validate:"required"`
}

type QuestionDetailRequest struct {
	CorrectAnswer string          `json:"correct_answer" validate:"required"`
	QuestionText  string          `json:"question_text" validate:"required"`
	Answer        []AnswerRequest `json:"answers" validate:"required"`
}
