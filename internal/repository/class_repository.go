package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type ClassRepository interface {
	Create(tx *gorm.DB, user *entity.Class) error
	Update(tx *gorm.DB, user *entity.Class) error
	Delete(tx *gorm.DB, user *entity.Class) error
}

type ClassRepositoryImpl struct {
	Repository[entity.Class]
}

func NewClassRepository() ClassRepository {
	return &ClassRepositoryImpl{}
}
