package model

type CourseResponse struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type CourseRequest struct {
	Name string `json:"name" validate:"required"`
}
