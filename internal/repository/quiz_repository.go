package repository

import (
	"strings"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type QuizRepository interface {
	Create(tx *gorm.DB, quiz *entity.Quiz) error
	Update(tx *gorm.DB, quiz *entity.Quiz) error
	Delete(tx *gorm.DB, quiz *entity.Quiz) error
	QuizDashboard(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint) error
	QuizDashboardActive(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint) error
	QuizDashboardArchived(tx *gorm.DB, quizzes *[]entity.Quiz, userId uint) error
}

type QuizRepositoryImpl struct {
	Repository[entity.Quiz]
}

func NewQuizRepository() QuizRepository {
	return &QuizRepositoryImpl{}
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
		Find(&quizzes).Error
}
