package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type LecturerTeachingRepository interface {
	CreateBacth(tx *gorm.DB, lecturerTeachings *[]entity.LecturerTeaching) error
	Delete(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error
	FindById(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error
	OneDataCheck(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error
	FindLecturerTeaching(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error
	FindLecturerClasses(tx *gorm.DB, userId uint, lecturerTeachings *[]entity.LecturerTeaching) error
	FindByCourseCode(tx *gorm.DB, userId uint, courseCode string) error
}

type LecturerTeachingRepositoryImpl struct {
	Repository[entity.LecturerTeaching]
}

func NewLecturerTeachingRepository() LecturerTeachingRepository {
	return &LecturerTeachingRepositoryImpl{}
}

// FindByCourseCode implements LecturerTeachingRepository.
func (repository *LecturerTeachingRepositoryImpl) FindByCourseCode(tx *gorm.DB, userId uint, courseCode string) error {
	return tx.First(&entity.LecturerTeaching{}, "user_id = ? AND course_code = ?", userId, courseCode).Error
}

// FindLecturerClasses implements LecturerTeachingRepository.
func (repository *LecturerTeachingRepositoryImpl) FindLecturerClasses(tx *gorm.DB, userId uint, lecturerTeachings *[]entity.LecturerTeaching) error {
	return tx.Where("user_id = ?", userId).Find(lecturerTeachings).Error
}

// CheckJoinData implements LecturerTeachingRepository.
func (repository *LecturerTeachingRepositoryImpl) FindById(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error {
	return tx.First(lecturerTeaching).Error
}

// OneDataCheck implements LecturerTeachingRepository.
func (repository *LecturerTeachingRepositoryImpl) OneDataCheck(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error {
	return tx.Where("(course_code = ? AND user_id = ? AND class_id = ?) OR (course_code = ? AND class_id = ?)",
		lecturerTeaching.CourseCode, lecturerTeaching.UserId, lecturerTeaching.ClassId,
		lecturerTeaching.CourseCode, lecturerTeaching.ClassId,
	).First(lecturerTeaching).Error
}

// FindLecturerTeaching implements LecturerTeachingRepository.
func (repository *LecturerTeachingRepositoryImpl) FindLecturerTeaching(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error {
	return tx.Where("course_code = ? AND user_id = ?", lecturerTeaching.CourseCode, lecturerTeaching.UserId).First(lecturerTeaching).Error
}
