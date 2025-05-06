package converter

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
)

func ClassToResponse(class *entity.Class) *model.ClassResponse {
	log.Println("log from class to response")

	return &model.ClassResponse{
		ID:   class.ID,
		Name: class.Name,
	}
}

func ClassToResponses(classs *[]entity.Class) *[]model.ClassResponse {
	var classResponses []model.ClassResponse

	log.Println("log from class to responses")

	for _, class := range *classs {
		classResponses = append(classResponses, *ClassToResponse(&class))
	}

	return &classResponses
}
