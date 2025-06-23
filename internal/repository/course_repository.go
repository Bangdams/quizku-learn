package repository

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(tx *gorm.DB, course *entity.Course) error
	Update(tx *gorm.DB, course *entity.Course) error
	Delete(tx *gorm.DB, course *entity.Course) error
	FindByCourseCode(tx *gorm.DB, course *entity.Course) error
	FindAll(tx *gorm.DB, courses *[]entity.Course) error
	FindAllByCourseCode(tx *gorm.DB, courseCode []string, courses *[]entity.Course) error
	FindByIdWithClass(tx *gorm.DB, courses *[]entity.Course, classId []uint) error
	FindWithClassSubject(tx *gorm.DB, count *int64, courseCode string, classId uint) error
	ListCoursesByUserWithClass(tx *gorm.DB, userId uint) ([]entity.Course, []uint, []uint, error)
	ListCoursesByUser(tx *gorm.DB, courses *[]entity.Course, userId uint) error
}

type CourseRepositoryImpl struct {
	Repository[entity.Course]
}

func NewCourseRepository() CourseRepository {
	return &CourseRepositoryImpl{}
}

// ListCoursesByUser implements CourseRepository.
func (repository *CourseRepositoryImpl) ListCoursesByUser(tx *gorm.DB, courses *[]entity.Course, userId uint) error {
	return tx.Table("courses").
		Select("courses.course_code, courses.name").
		Joins("JOIN class_subjects ON courses.course_code = class_subjects.course_code").
		Joins("JOIN user_classes ON class_subjects.class_id = user_classes.class_id").
		Where("user_classes.user_id = ?", userId).
		Find(&courses).Error
}

// FindWithClassSubject implements CourseRepository.
func (repository *CourseRepositoryImpl) FindWithClassSubject(tx *gorm.DB, count *int64, courseCode string, classId uint) error {
	tx.Table("class_subjects").Where("course_code = ?", courseCode).
		Where("class_id = ?", classId).Count(count)

	return nil
}

// FindByIdWithClass implements ClassRepository.
func (repository *CourseRepositoryImpl) FindByIdWithClass(tx *gorm.DB, courses *[]entity.Course, classId []uint) error {
	return tx.Preload("Classes", "id IN ?", classId).Find(courses).Error
	// return tx.Preload("Classes", "id NOT IN (SELECT id FROM lecturer_teachings WHERE class_id IN ?)", classId).
	// 	Find(&courses).Error
}

// ListCoursesByUserWithClass implements ClassRepository.
func (repository *CourseRepositoryImpl) ListCoursesByUserWithClass(tx *gorm.DB, userId uint) ([]entity.Course, []uint, []uint, error) {
	var lecturerTeachings []entity.LecturerTeaching
	err := tx.Preload("Class").Preload("Course").
		Where("user_id = ?", userId).
		Find(&lecturerTeachings).Error

	if err != nil {
		return nil, nil, nil, err
	}

	courseMap := make(map[string]*entity.Course)
	for _, lt := range lecturerTeachings {
		if _, exists := courseMap[lt.Course.CourseCode]; !exists {
			courseMap[lt.Course.CourseCode] = &lt.Course
		}

		courseMap[lt.Course.CourseCode].Classes = append(courseMap[lt.Course.CourseCode].Classes, lt.Class)
	}

	var studentCount int64
	var totalStudents []uint
	var quizCount int64
	var totalQuiz []uint
	var courses []entity.Course

	for _, course := range courseMap {
		log.Println(course)
		// student count
		for _, class := range course.Classes {
			var count int64
			tx.Model(&entity.UserClass{}).Where("class_id = ?", class.ID).Count(&count)
			studentCount += count

			var qCount int64
			tx.Model(&entity.Quiz{}).Where("class_id = ?", class.ID).Where("course_code = ?", course.CourseCode).Count(&qCount)
			quizCount += qCount
		}

		totalStudents = append(totalStudents, uint(studentCount))
		totalQuiz = append(totalQuiz, uint(quizCount))
		courses = append(courses, *course)

		studentCount = 0
		quizCount = 0
	}

	return courses, totalStudents, totalQuiz, nil
}

// FindAllByCourseCode implements CourseRepository.
func (repository *CourseRepositoryImpl) FindAllByCourseCode(tx *gorm.DB, courseCode []string, courses *[]entity.Course) error {
	return tx.Model(&entity.Course{}).Where("course_code IN ?", courseCode).Find(courses).Error
}

// FindAll implements CourseRepository.
func (repository *CourseRepositoryImpl) FindAll(tx *gorm.DB, courses *[]entity.Course) error {
	return tx.Find(courses).Error
}

// FindByCourseCode implements CourseRepository.
func (repository *CourseRepositoryImpl) FindByCourseCode(tx *gorm.DB, course *entity.Course) error {
	return tx.First(course, "course_code=?", course.CourseCode).Error
}
