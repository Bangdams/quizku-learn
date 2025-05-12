package model

type AnswerResponse struct {
	ID               uint   `json:"id" validate:"required"`
	QuestionDetailId uint   `json:"question_detail_id" validate:"required"`
	Answer           string `json:"answer" validate:"required"`
}

type AnswerRequest struct {
	Choice string `json:"choice" validate:"required"`
	Answer string `json:"answer_text" validate:"required"`
}
