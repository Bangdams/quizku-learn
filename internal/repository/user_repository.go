package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(tx *gorm.DB, user *entity.User) error
	Update(tx *gorm.DB, user *entity.User) error
	Delete(tx *gorm.DB, user *entity.User) error
	Login(tx *gorm.DB, user *entity.User, keyword string) error
	FindByEmail(tx *gorm.DB, user *entity.User) error
	FindAll(tx *gorm.DB, userId uint, users *[]entity.User) error
	FindByRole(tx *gorm.DB, role string, userId uint, users *[]entity.User) error
	Search(tx *gorm.DB, users *[]entity.User, name string) error
	FindById(tx *gorm.DB, user *entity.User) error
	AdminDashboardReport(tx *gorm.DB, response *model.AdminDashboardReportResponse)
	LecturerDashboardReport(tx *gorm.DB, userId uint, classIds []uint, response *model.LecturerDashboardReportResponse)
}

type UserRepositoryImpl struct {
	Repository[entity.User]
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

// AdminDashboardReport implements UserRepository.
func (repository *UserRepositoryImpl) AdminDashboardReport(tx *gorm.DB, response *model.AdminDashboardReportResponse) {
	tx.Model(&entity.User{}).Count(&response.TotalUsers)
	tx.Model(&entity.Quiz{}).Count(&response.TotalQuizzes)
	tx.Model(&entity.Class{}).Count(&response.TotalClasses)
	tx.Model(&entity.Question{}).Count(&response.TotalQuestions)
}

// LecturerDashboardReport implements UserRepository.
func (repository *UserRepositoryImpl) LecturerDashboardReport(tx *gorm.DB, userId uint, classIds []uint, response *model.LecturerDashboardReportResponse) {
	tx.Model(&entity.User{}).
		Joins("JOIN user_classes ON users.id = user_classes.user_id").
		Where("user_classes.id IN ?", classIds).Count(&response.TotalUsers)

	tx.Model(&entity.Quiz{}).
		Count(&response.TotalQuizzes)

	tx.Model(&entity.LecturerTeaching{}).
		Where("user_id = ?", userId).
		Group("class_id").
		Count(&response.TotalClasses)

	tx.Model(&entity.LecturerTeaching{}).
		Where("user_id = ?", userId).
		Select("COUNT(DISTINCT course_code)").
		Count(&response.TotalCourses)
}

// FindByRole implements UserRepository.
func (repository *UserRepositoryImpl) FindByRole(tx *gorm.DB, role string, userId uint, users *[]entity.User) error {
	return tx.Where("role = ?", role).Not("id = ?", userId).Find(users).Error
}

// FindAll implements UserRepository.
func (repository *UserRepositoryImpl) FindAll(tx *gorm.DB, userId uint, users *[]entity.User) error {
	return tx.Not("id = ?", userId).Find(users).Error
}

// FindById implements UserRepository.
func (repository *UserRepositoryImpl) FindById(tx *gorm.DB, user *entity.User) error {
	return tx.First(user).Error
}

// Login implements UserRepository.
func (*UserRepositoryImpl) Login(tx *gorm.DB, user *entity.User, keyword string) error {
	return tx.Where("email = ?", keyword).First(user).Error
}

// FindByEmail implements UserRepository.
func (repository *UserRepositoryImpl) FindByEmail(tx *gorm.DB, user *entity.User) error {
	return tx.First(user, "email=?", user.Email).Error
}

// Search implements UserRepository.
func (*UserRepositoryImpl) Search(tx *gorm.DB, users *[]entity.User, name string) error {
	return tx.Where("name LIKE ?", "%"+name+"%").Find(&users).Error
}
