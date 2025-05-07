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

func ClassSubjectToResponse(class *entity.Class, courses *[]entity.Course) *model.ClassSubjectResponse {
	log.Println("log from class to response")

	var courseCodes []string
	for _, v := range *courses {
		courseCodes = append(courseCodes, v.CourseCode)
	}

	return &model.ClassSubjectResponse{
		ClassId:     int(class.ID),
		CourseCodes: courseCodes,
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
