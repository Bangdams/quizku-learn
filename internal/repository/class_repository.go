package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type ClassRepository interface {
	Create(tx *gorm.DB, class *entity.Class) error
	Update(tx *gorm.DB, class *entity.Class) error
	Delete(tx *gorm.DB, class *entity.Class) error
	FindAll(tx *gorm.DB, classes *[]entity.Class) error
	FindByName(tx *gorm.DB, class *entity.Class) error
	FindById(tx *gorm.DB, class *entity.Class) error
}

type ClassRepositoryImpl struct {
	Repository[entity.Class]
}

func NewClassRepository() ClassRepository {
	return &ClassRepositoryImpl{}
}

// FindByClassCode implements ClassRepository.
func (repository *ClassRepositoryImpl) FindAll(tx *gorm.DB, classes *[]entity.Class) error {
	return tx.Find(classes).Error
}

// FindByName implements ClassRepository.
func (repository *ClassRepositoryImpl) FindByName(tx *gorm.DB, class *entity.Class) error {
	return tx.First(class, "name=?", class.Name).Error
}

// FindById implement ClassRepository
func (repository *ClassRepositoryImpl) FindById(tx *gorm.DB, class *entity.Class) error {
	return tx.First(class).Error
}
