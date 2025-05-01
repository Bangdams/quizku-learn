package model

type QuestionResponse struct {
	ID            uint `json:"id" validate:"required"`
	Name          uint `json:"name" validate:"required"`
	QuestionCount uint `json:"question_count" validate:"required"`
	Duration      uint `json:"duration" validate:"required"`
	CourseId      uint `json:"course_id" validate:"required"`
	UserId        uint `json:"user_id" validate:"required"`
}

type QuestionRequest struct {
	ID            uint `json:"id" validate:"required"`
	Name          uint `json:"name" validate:"required"`
	QuestionCount uint `json:"question_count" validate:"required"`
	Duration      uint `json:"duration" validate:"required"`
	CourseId      uint `json:"course_id" validate:"required"`
	UserId        uint `json:"user_id" validate:"required"`
}
