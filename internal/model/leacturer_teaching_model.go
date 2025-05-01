package model

type LecturerTeachingResponse struct {
	ID      uint   `json:"id" validate:"required"`
	Code    string `json:"code" validate:"required"`
	UserId  uint   `json:"user_id" validate:"required"`
	ClassId uint   `json:"class_id" validate:"required"`
}

type LecturerTeachingRequest struct {
	Code    string `json:"code" validate:"required"`
	UserId  uint   `json:"user_id" validate:"required"`
	ClassId uint   `json:"class_id" validate:"required"`
}
