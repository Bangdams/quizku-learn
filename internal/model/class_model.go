package model

type ClassResponse struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type ClassRequest struct {
	Name string `json:"name" validate:"required"`
}
