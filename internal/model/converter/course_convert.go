package converter

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
)

func CourseToResponse(course *entity.Course) *model.CourseResponse {
	log.Println("log from course to response")

	return &model.CourseResponse{
		CourseCode: course.CourseCode,
		Name:       course.Name,
	}
}

func CourseToResponses(courses *[]entity.Course) *[]model.CourseResponse {
	var courseResponses []model.CourseResponse

	log.Println("log from course to responses")

	for _, course := range *courses {
		courseResponses = append(courseResponses, *CourseToResponse(&course))
	}

	return &courseResponses
}
