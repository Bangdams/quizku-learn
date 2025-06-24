package repository

import (
	"strings"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"gorm.io/gorm"
)

type QuizRepository interface {
	Create(tx *gorm.DB, quiz *entity.Quiz) error
	Update(tx *gorm.DB, quiz *entity.Quiz) error
	Delete(tx *gorm.DB, quiz *entity.Quiz) error
	FindById(tx *gorm.DB, quiz *entity.Quiz) error
	QuizDashboard(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint) error
	QuizDashboardActive(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint) error
	QuizDashboardArchived(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint) error
	FindByUserAndCourse(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint, courseCode string) error
	FindAll(tx *gorm.DB, quizzes *[]entity.Quiz) error
	StartQuiz(tx *gorm.DB, response *model.StartQuizResponse, quizId uint) error
}

type QuizRepositoryImpl struct {
	Repository[entity.Quiz]
}

func NewQuizRepository() QuizRepository {
	return &QuizRepositoryImpl{}
}

// StartQuiz implements QuizRepository.
func (repository *QuizRepositoryImpl) StartQuiz(tx *gorm.DB, response *model.StartQuizResponse, quizId uint) error {
	quiz := &entity.Quiz{}
	err := tx.Preload("Course").
		Preload("Question.User").
		Where("quizzes.id = ?", quizId).
		First(quiz).Error
	if err != nil {
		return err
	}

	response.CourseName = quiz.Course.Name
	response.CourseCode = quiz.Course.CourseCode
	response.QuizName = quiz.Question.Name
	response.LecturerName = quiz.Question.User.Name
	response.Duration = quiz.Question.Duration
	response.QuestionCount = quiz.Question.QuestionCount

	questionDetails := &[]entity.QuestionDetail{}
	err = tx.
		Preload("Answers").
		Where("question_details.question_id = ?", quiz.Question.ID).
		Find(questionDetails).Error
	if err != nil {
		return err
	}

	for _, data := range *questionDetails {
		choiceItems := []model.ChoiceItem{}

		for _, choiceItem := range data.Answers {
			choiceItems = append(choiceItems, model.ChoiceItem{
				AnswerId: choiceItem.ID,
				Choice:   choiceItem.Choice,
				Answer:   choiceItem.Answer,
			})
		}

		response.QuestionItems = append(response.QuestionItems, model.QuestionItem{
			QuestionText: data.QuestionText,
			ChoiceItems:  choiceItems,
		})
	}

	return nil
}

// FindAll implements QuizRepository.
func (repository *QuizRepositoryImpl) FindAll(tx *gorm.DB, quizzes *[]entity.Quiz) error {
	return tx.Preload("Course").
		Preload("Question").
		Preload("Class.UserClasses").
		Find(quizzes).Error
}

// FindByUserAndCourse implements QuizRepository.
func (repository *QuizRepositoryImpl) FindByUserAndCourse(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint, courseCode string) error {
	return tx.
		Joins("JOIN user_classes ON user_classes.class_id = quizzes.class_id").
		Joins("JOIN courses ON quizzes.course_code = courses.course_code").
		Joins("JOIN questions ON quizzes.question_id = questions.id").
		Where("user_classes.user_id = ? AND courses.course_code = ?", userId, courseCode).
		Where("quizzes.deadline > NOW()").
		Where(`
        NOT EXISTS (
            SELECT 1 FROM quizz_results
            WHERE quizz_results.quizz_id = quizzes.id AND quizz_results.user_id = ?
        )
    `, userId).
		Preload("Class.UserClasses").
		Preload("Course").
		Preload("Question.User").
		Find(&quizzes).Error

}

// FindById implements QuizRepository.
func (repository *QuizRepositoryImpl) FindById(tx *gorm.DB, quiz *entity.Quiz) error {
	return tx.Preload("Question.QuestionDetails").First(quiz).Error
}

// QuizDashboardActive implements QuizRepository.
func (repository *QuizRepositoryImpl) QuizDashboardActive(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint) error {
	lecturerTeachings := &[]entity.LecturerTeaching{}
	uniqueCourseCodes := make(map[string]bool)

	var courseCodes []string
	var classIds []uint

	tx.Preload("Course").
		Where("user_id = ?", userId).
		Group("course_code").
		Find(&lecturerTeachings)

	for _, value := range *lecturerTeachings {
		if !uniqueCourseCodes[value.CourseCode] {
			uniqueCourseCodes[value.CourseCode] = true
			courseCodes = append(courseCodes, value.CourseCode)
		}

		classIds = append(classIds, value.ClassId)
	}

	var conditions []string
	var values []interface{}

	for i := range courseCodes {
		conditions = append(conditions, "(course_code = ? AND class_id = ?)")
		values = append(values, courseCodes[i], classIds[i])
	}

	query := strings.Join(conditions, " OR ")

	return tx.Preload("Course").
		Preload("Question").
		Preload("Class.UserClasses").
		Where(query, values...).
		Where("deadline >= CURRENT_DATE").
		Order("deadline desc").
		Find(&quizzes).Error
}

// QuizDashboardArchived implements QuizRepository.
func (repository *QuizRepositoryImpl) QuizDashboardArchived(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint) error {
	lecturerTeachings := &[]entity.LecturerTeaching{}
	uniqueCourseCodes := make(map[string]bool)

	var courseCodes []string
	var classIds []uint

	tx.Preload("Course").
		Where("user_id = ?", userId).
		Group("course_code").
		Find(&lecturerTeachings)

	for _, value := range *lecturerTeachings {
		if !uniqueCourseCodes[value.CourseCode] {
			uniqueCourseCodes[value.CourseCode] = true
			courseCodes = append(courseCodes, value.CourseCode)
		}

		classIds = append(classIds, value.ClassId)
	}

	var conditions []string
	var values []interface{}

	for i := range courseCodes {
		conditions = append(conditions, "(course_code = ? AND class_id = ?)")
		values = append(values, courseCodes[i], classIds[i])
	}

	query := strings.Join(conditions, " OR ")

	return tx.Preload("Course").
		Preload("Question").
		Preload("Class.UserClasses").
		Where(query, values...).
		Where("deadline < CURRENT_DATE").
		Order("deadline desc").
		Find(&quizzes).Error
}

// QuizDashboard implements QuizRepository.
func (repository *QuizRepositoryImpl) QuizDashboard(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint) error {
	lecturerTeachings := &[]entity.LecturerTeaching{}
	uniqueCourseCodes := make(map[string]bool)

	var courseCodes []string
	var classIds []uint

	tx.Preload("Course").
		Where("user_id = ?", userId).
		Group("course_code").
		Find(&lecturerTeachings)

	for _, value := range *lecturerTeachings {
		if !uniqueCourseCodes[value.CourseCode] {
			uniqueCourseCodes[value.CourseCode] = true
			courseCodes = append(courseCodes, value.CourseCode)
		}

		classIds = append(classIds, value.ClassId)
	}

	var conditions []string
	var values []interface{}

	for i := range courseCodes {
		conditions = append(conditions, "(course_code = ? AND class_id = ?)")
		values = append(values, courseCodes[i], classIds[i])
	}

	query := strings.Join(conditions, " OR ")

	return tx.Preload("Course").
		Preload("Question").
		Preload("Class.UserClasses").
		Where(query, values...).
		Order("deadline desc").
		Find(&quizzes).Error
}
