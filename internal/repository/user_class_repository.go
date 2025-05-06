package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type UserClassRepository interface {
	Create(tx *gorm.DB, userClass *entity.UserClass) error
}

type UserClassRepositoryImpl struct {
	Repository[entity.UserClass]
}

func NewUserClassRepository() UserClassRepository {
	return &UserClassRepositoryImpl{}
}
