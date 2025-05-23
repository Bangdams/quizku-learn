package repository

import (
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
	ListCoursesByUser(tx *gorm.DB, userId uint) ([]entity.Course, []uint, error)
}

type CourseRepositoryImpl struct {
	Repository[entity.Course]
}

func NewCourseRepository() CourseRepository {
	return &CourseRepositoryImpl{}
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
}

// ListCoursesByUser implements ClassRepository.
func (repository *CourseRepositoryImpl) ListCoursesByUser(tx *gorm.DB, userId uint) ([]entity.Course, []uint, error) {
	var lecturerTeachings []entity.LecturerTeaching
	err := tx.Preload("Class").Preload("Course").
		Where("user_id = ?", userId).
		Find(&lecturerTeachings).Error

	if err != nil {
		return nil, nil, err
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
	var courses []entity.Course

	for _, course := range courseMap {
		// student count
		for _, class := range course.Classes {
			var count int64
			tx.Model(&entity.UserClass{}).Where("class_id = ?", class.ID).Count(&count)
			studentCount += count
		}

		totalStudents = append(totalStudents, uint(studentCount))
		courses = append(courses, *course)

		studentCount = 0
	}

	return courses, totalStudents, nil
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
