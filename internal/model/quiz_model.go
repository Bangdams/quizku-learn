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

type QuizHistoryStudentResponse struct {
	ID         uint   `json:"id" validate:"required"`
	QuizName   string `json:"quiz_name" validate:"required"`
	CourseName string `json:"course_name" validate:"required"`
	CreatedAt  string `json:"created_at" validate:"required"`
	Score      uint   `json:"score" validate:"required"`
	Status     string `json:"status" validate:"required"`
}

type QuizHistoryResponse struct {
	ID            uint      `json:"id" validate:"required"`
	CourseName    string    `json:"course_name" validate:"required"`
	QuizName      string    `json:"quiz_name" validate:"required"`
	StudentCount  uint      `json:"student_count" validate:"required"`
	QuestionCount uint      `json:"question_count" validate:"required"`
	CreatedAt     time.Time `json:"created_at" validate:"required"`
}

type QuizStudentResultResponse struct {
	Score                uint               `json:"score" validate:"required"`
	QuestionCount        uint               `json:"question_count" validate:"required"`
	CountCorrectAnswer   uint               `json:"count_correct_answer" validate:"required"`
	CountIncorrectAnswer uint               `json:"count_incorrect_answer" validate:"required"`
	AnswerUserDetails    []AnswerUserDetail `json:"incorrect_answer_detail" validate:"required"`
}

type ScoreResult struct {
	Score int
	Type  string
}

type QuizResultAnalysisResponse struct {
	AvarageScore        uint                        `json:"avarage_score" validate:"required"`
	HighestScore        uint                        `json:"highest_score" validate:"required"`
	LowestScore         uint                        `json:"lowest_score" validate:"required"`
	CompletionRate      uint                        `json:"completion_rate" validate:"required"`
	StudentCount        uint                        `json:"student_count" validate:"required"`
	UserCompleteCount   uint                        `json:"user_complete_count" validate:"required"`
	UserIncompleteCount uint                        `json:"user_incomplete_count" validate:"required"`
	CoursesName         string                      `json:"course_name" validate:"required"`
	QuizName            string                      `json:"quiz_name" validate:"required"`
	QuestionCount       uint                        `json:"question_count" validate:"required"`
	CreatedAt           string                      `json:"created_at" validate:"required"`
	Duration            uint                        `json:"duration" validate:"required"`
	Students            []QuizResultAnalysisStudent `json:"students" validate:"required"`
}

type QuizResultAnalysisStudent struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	ClassName string `json:"class_name" validate:"required"`
	Score     uint   `json:"score" validate:"required"`
	Status    string `json:"status" validate:"required"`
}

type AnswerUserDetail struct {
	QuestionText    string `json:"question_text" validate:"required"`
	CorrectAnswer   string `json:"correct_answer" validate:"required"`
	IncorrectAnswer string `json:"incorrect_answer" validate:"required"`
}

type QuizStudentResponse struct {
	ID            uint   `json:"id" validate:"required"`
	LecturerName  string `json:"lecturer_name" validate:"required"`
	CourseName    string `json:"course_name" validate:"required"`
	QuestionName  string `json:"question_name" validate:"required"`
	Deadline      string `json:"deadline" validate:"required"`
	QuestionCount uint   `json:"question_count" validate:"required"`
	Duration      uint   `json:"duration" validate:"required"`
}

type QuizDashboardResponse struct {
	ID            uint      `json:"id" validate:"required"`
	QuestionName  string    `json:"question_name" validate:"required"`
	CourseName    string    `json:"course_name" validate:"required"`
	QuestionCount uint      `json:"question_count" validate:"required"`
	StudentCount  int       `json:"student_count" validate:"required"`
	Deadline      time.Time `json:"deadline" validate:"required"`
	CreatedAt     time.Time `json:"created_at" validate:"required"`
}

type QuizRequest struct {
	CourseCode string `json:"course_code" validate:"required"`
	ClassId    uint   `json:"class_id" validate:"required"`
	QuestionId uint   `json:"question_id" validate:"required"`
	Deadline   string `json:"deadline" validate:"required"`
}

type StartQuizResponse struct {
	CourseName    string         `json:"course_name" validate:"required"`
	CourseCode    string         `json:"course_code" validate:"required"`
	QuizName      string         `json:"quiz_name" validate:"required"`
	LecturerName  string         `json:"lecturer_name" validate:"required"`
	Duration      uint           `json:"duration" validate:"required"`
	QuestionCount uint           `json:"question_count" validate:"required"`
	QuestionItems []QuestionItem `json:"question_items" validate:"required"`
}

type QuestionItem struct {
	QuestionText string       `json:"question_text" validate:"required"`
	ChoiceItems  []ChoiceItem `json:"choice_items" validate:"required"`
}

type ChoiceItem struct {
	Choice string `json:"choice" validate:"required"`
	Answer string `json:"answer" validate:"required"`
}
