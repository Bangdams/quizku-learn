package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type LecturerTeachingRepository interface {
	Create(tx *gorm.DB, user *entity.LecturerTeaching) error
	Update(tx *gorm.DB, user *entity.LecturerTeaching) error
	Delete(tx *gorm.DB, user *entity.LecturerTeaching) error
}

type LecturerTeachingRepositoryImpl struct {
	Repository[entity.LecturerTeaching]
}

func NewLecturerTeachingRepository() LecturerTeachingRepository {
	return &LecturerTeachingRepositoryImpl{}
}
