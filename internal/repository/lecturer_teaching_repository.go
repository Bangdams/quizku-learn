package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type LecturerTeachingRepository interface {
	CreateBacth(tx *gorm.DB, lecturerTeachings *[]entity.LecturerTeaching) error
	Delete(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error
	OneDataCheck(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error
	FindById(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error
}

type LecturerTeachingRepositoryImpl struct {
	Repository[entity.LecturerTeaching]
}

func NewLecturerTeachingRepository() LecturerTeachingRepository {
	return &LecturerTeachingRepositoryImpl{}
}

// CheckJoinData implements LecturerTeachingRepository.
func (repository *LecturerTeachingRepositoryImpl) FindById(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error {
	return tx.First(lecturerTeaching).Error
}

// OneDataCheck implements LecturerTeachingRepository.
func (repository *LecturerTeachingRepositoryImpl) OneDataCheck(tx *gorm.DB, lecturerTeaching *entity.LecturerTeaching) error {
	return tx.First(lecturerTeaching, "course_code = ? AND user_id = ? AND class_id = ?", lecturerTeaching.CourseCode, lecturerTeaching.UserId, lecturerTeaching.ClassId).Error
}
