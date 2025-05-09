package converter

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	log.Println("log from user to response")

	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Image:     user.Image,
		CreatedAt: user.CreatedAt.Format("2006-01-02"),
	}
}

func LoginUserToResponse(token string) *model.LoginResponse {
	log.Println("log from login user to response")

	return &model.LoginResponse{
		AccessToken: token,
	}
}

func UserToResponses(users *[]entity.User) *[]model.UserResponse {
	var userResponses []model.UserResponse

	log.Println("log from user to responses")

	for _, user := range *users {
		userResponses = append(userResponses, *UserToResponse(&user))
	}

	return &userResponses
}
