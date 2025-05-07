package model

type ClassResponse struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type ClassRequest struct {
	Name string `json:"name" validate:"required"`
}

type ClassUpdate struct {
	ClassResponse
}

type ClassSubjectResponse struct {
	ClassId     int      `json:"class_id" validate:"required"`
	CourseCodes []string `json:"course_code" validate:"required"`
}

type ClassSubjectRequest struct {
	ClassId     int      `json:"class_id" validate:"required"`
	CourseCodes []string `json:"course_code" validate:"required"`
}
