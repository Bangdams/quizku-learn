package model

type UserAnswerResponse struct {
	ID       uint `json:"id" validate:"required"`
	AnswerId uint `json:"answer_id" validate:"required"`
	QuizzId  uint `json:"quiz_id" validate:"required"`
}

type UserAnswerRequests struct {
	UserAnswers []UserAnswerRequest `json:"user_answers" validate:"required"`
}

type UserAnswerRequest struct {
	AnswerId uint `json:"answer_id" validate:"required"`
}
