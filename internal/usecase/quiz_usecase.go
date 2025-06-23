package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/model/converter"
	"github.com/Bangdams/quizku-learn/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type QuizUsecase interface {
	Create(ctx context.Context, request *model.QuizRequest, userId uint) (*model.QuizResponse, error)
	Update(ctx context.Context, request *model.QuizRequest) (*model.QuizResponse, error)
	Delete(ctx context.Context, quizId uint) error
	QuizDashboard(ctx context.Context, userId uint, query string) (*[]model.QuizDashboardResponse, error)
	FindByUserAndCourse(ctx context.Context, userId uint, courseCode string) (*[]model.QuizStudentResponse, error)
	QuizStudentResult(ctx context.Context, userId uint, quizId uint) (*model.QuizStudentResultResponse, error)
	QuizResultAnalysis(ctx context.Context, quizId uint) (*model.QuizResultAnalysisResponse, error)
	FindAll(ctx context.Context) (*[]model.QuizHistoryResponse, error)
	GetQuizHistoryForStudent(ctx context.Context, userId uint) (*[]model.QuizHistoryStudentResponse, error)
	StartQuiz(ctx context.Context, quizId uint) (*model.StartQuizResponse, error)
}

type QuizUsecaseImpl struct {
	QuizResultRepo       repository.QuizResultRepository
	AnswerRepo           repository.AnswerRepository
	UserAnswerRepo       repository.UserAnswerRepository
	QuizRepo             repository.QuizRepository
	ClassRepo            repository.ClassRepository
	QuestionRepo         repository.QuestionRepository
	LecturerTeachingRepo repository.LecturerTeachingRepository
	DB                   *gorm.DB
	Validate             *validator.Validate
}

func NewQuizUsecase(quizResultRepo repository.QuizResultRepository, answerRepo repository.AnswerRepository, userAnswerRepo repository.UserAnswerRepository, quizRepo repository.QuizRepository, classRepo repository.ClassRepository, questionRepo repository.QuestionRepository, LecturerTeachingRepo repository.LecturerTeachingRepository, DB *gorm.DB, validate *validator.Validate) QuizUsecase {
	return &QuizUsecaseImpl{
		QuizResultRepo:       quizResultRepo,
		AnswerRepo:           answerRepo,
		UserAnswerRepo:       userAnswerRepo,
		QuizRepo:             quizRepo,
		ClassRepo:            classRepo,
		QuestionRepo:         questionRepo,
		LecturerTeachingRepo: LecturerTeachingRepo,
		DB:                   DB,
		Validate:             validate,
	}
}

// StartQuiz implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) StartQuiz(ctx context.Context, quizId uint) (*model.StartQuizResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	startQuizResponse := &model.StartQuizResponse{}

	err := quizUsecase.QuizRepo.StartQuiz(tx, startQuizResponse, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "quiz data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id quiz : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by id quiz : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return startQuizResponse, nil
}

// GetQuizHistoryForStudent implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) GetQuizHistoryForStudent(ctx context.Context, userId uint) (*[]model.QuizHistoryStudentResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	quizResults := &[]entity.QuizzResult{}

	err := quizUsecase.QuizResultRepo.GetQuizHistoryForStudent(tx, userId, quizResults)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "quiz data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id quiz : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by id quiz : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success show GetQuizHistoryForStudent from usecase quiz")
	return converter.QuizHistoryStudentToResponses(quizResults), nil
}

// FindAll implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) FindAll(ctx context.Context) (*[]model.QuizHistoryResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	quizzes := &[]entity.Quiz{}

	err := quizUsecase.QuizRepo.FindAll(tx, quizzes)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "quiz data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id quiz : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by id quiz : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success show QuizResultAnalysis from usecase quiz")
	return converter.QuizHistoryToResponses(quizzes), nil
}

// QuizResultAnalysis implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) QuizResultAnalysis(ctx context.Context, quizId uint) (*model.QuizResultAnalysisResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	response, err := quizUsecase.QuizResultRepo.QuizResultAnalysis(tx, quizId, &[]entity.QuizzResult{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "quiz data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id quiz : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by id quiz : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success show QuizResultAnalysis from usecase quiz")
	return &response, nil
}

// QuizStudentResult implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) QuizStudentResult(ctx context.Context, userId uint, quizId uint) (*model.QuizStudentResultResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	quiz := &entity.Quiz{
		ID: quizId,
	}

	err := quizUsecase.QuizRepo.FindById(tx, quiz)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "quiz data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id quiz : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by id quiz : ", err)
		return nil, fiber.ErrInternalServerError
	}

	err = quizUsecase.UserAnswerRepo.VerifyUserClassQuiz(tx, userId, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "You do not have access to this quiz because it is not assigned to your class."
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error verify user class : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error verify user class : ", err)
		return nil, fiber.ErrInternalServerError
	}

	userAnswers := []entity.UserAnswer{}
	err = quizUsecase.UserAnswerRepo.GetUserAnswersByQuestion(tx, &userAnswers, quiz.QuestionId, userId)
	if err != nil {
		log.Println("failed when GetUserAnswersByQuestion : ", err)
		return nil, fiber.ErrInternalServerError
	}

	for i, element := range userAnswers {
		questionDetailsId := element.Answer.QuestionDetail.ID
		correctChoice := element.Answer.QuestionDetail.CorrectAnswer

		answer := entity.Answer{}
		err = quizUsecase.AnswerRepo.GetAnswersByQuestionDetail(tx, questionDetailsId, correctChoice, &answer)
		if err != nil {
			log.Println("failed when GetAnswersByQuestionDetail : ", err)
			return nil, fiber.ErrInternalServerError
		}

		userAnswers[i].Answer.QuestionDetail.Answers = append(userAnswers[i].Answer.QuestionDetail.Answers, answer)
	}

	quizResult := entity.QuizzResult{}
	err = quizUsecase.QuizResultRepo.FindByUserAndQuiz(tx, userId, quizId, &quizResult)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "quizz result data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by user and quiz : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by user and quiz : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success show quiz strudent result from usecase quiz")
	return converter.QuizStudentResultToResponse(&userAnswers, &quizResult), nil
}

// FindByUserAndCourse implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) FindByUserAndCourse(ctx context.Context, userId uint, courseCode string) (*[]model.QuizStudentResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var quizzes = &[]entity.Quiz{}

	err := quizUsecase.QuizRepo.FindByUserAndCourse(tx, quizzes, userId, courseCode)
	if err != nil {
		log.Println("failed when find by user and course : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find by user and course from usecase quiz")
	return converter.QuizStudentToResponses(quizzes), nil
}

// QuizDashboard implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) QuizDashboard(ctx context.Context, userId uint, query string) (*[]model.QuizDashboardResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var quizzes = &[]entity.Quiz{}

	if query == "active" {
		err := quizUsecase.QuizRepo.QuizDashboardActive(tx, quizzes, userId)
		if err != nil {
			log.Println("failed when find all repo quiz : ", err)
			return nil, fiber.ErrInternalServerError
		}
	} else if query == "archived" {
		err := quizUsecase.QuizRepo.QuizDashboardArchived(tx, quizzes, userId)
		if err != nil {
			log.Println("failed when find all repo quiz : ", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		err := quizUsecase.QuizRepo.QuizDashboard(tx, quizzes, userId)
		if err != nil {
			log.Println("failed when find all repo quiz : ", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.QuizDashboardResponses(quizzes), nil
}

// Create implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) Create(ctx context.Context, request *model.QuizRequest, userId uint) (*model.QuizResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := quizUsecase.Validate.Struct(request)
	if err != nil {
		log.Println("Invalid request Body : ", err)
		return nil, fiber.ErrBadRequest
	}

	lecturerTeachings := &[]entity.LecturerTeaching{}

	var errorResponse model.ErrorResponse

	err = quizUsecase.LecturerTeachingRepo.FindLecturerTeaching(tx, request.CourseCode, userId, lecturerTeachings)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "Data not found"
			errorResponse.Details = []string{
				fmt.Sprintf("The lecturer is not assigned to teach the course: %s.", request.CourseCode),
			}

			jsonString, _ := json.Marshal(errorResponse)
			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error FindLecturerTeaching from question usecase : ", err)
		return nil, fiber.ErrInternalServerError
	}

	class := &entity.Class{
		ID: request.ClassId,
	}

	err = quizUsecase.ClassRepo.FindById(tx, class)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "Class data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)
			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by class from quiz usecase : ", err)
		return nil, fiber.ErrInternalServerError
	}

	found := false
	for _, lecturerTeaching := range *lecturerTeachings {
		if lecturerTeaching.ClassId == request.ClassId {
			found = true
			break
		}
	}

	if !found {
		errorResponse.Message = "Data not found"
		errorResponse.Details = []string{
			fmt.Sprintf("The lecturer does not teach in this class : %s", class.Name),
		}

		jsonString, _ := json.Marshal(errorResponse)
		return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	question := &entity.Question{
		ID:         request.QuestionId,
		CourseCode: request.CourseCode,
	}

	err = quizUsecase.QuestionRepo.FindByCourseAndId(tx, question)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "question data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)
			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by question from quiz usecase : ", err)
		return nil, fiber.ErrInternalServerError
	}

	// Format yang sesuai dengan input string
	layout := "2006-01-02"

	// Parsing string menjadi time.Time
	parsedDate, err := time.Parse(layout, request.Deadline)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return nil, fiber.NewError(fiber.ErrBadRequest.Code, "bad request")
	}

	quiz := &entity.Quiz{
		CourseCode: request.CourseCode,
		ClassId:    request.ClassId,
		QuestionId: request.QuestionId,
		Deadline:   parsedDate,
		Status:     "aktif",
	}

	err = quizUsecase.QuizRepo.Create(tx, quiz)
	if err != nil {
		log.Println("failed when create repo quiz : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase quiz")

	return converter.QuizToResponse(quiz), nil
}

// Delete implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) Delete(ctx context.Context, quizId uint) error {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	quiz := &entity.Quiz{}
	quiz.ID = quizId

	err := quizUsecase.QuizRepo.FindById(tx, quiz)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Quiz data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error delete quiz : ", err)

			return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by id from quiz usecase : ", err)
		return fiber.ErrInternalServerError
	}

	err = quizUsecase.QuizRepo.Delete(tx, quiz)
	if err != nil {
		log.Println("failed when delete repo quiz : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase quiz")

	return nil
}

// Update implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) Update(ctx context.Context, request *model.QuizRequest) (*model.QuizResponse, error) {
	panic("unimplemented")
}
