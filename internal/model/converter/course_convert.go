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

func UserCourseListToResponse(courses *[]entity.Course, totalStudents *[]uint) *model.UserCourseListResponse {
	var dataItem model.UserCourseListItem
	var dataItems []model.UserCourseListItem

	log.Println("log from course to response")

	for i, course := range *courses {
		dataItem.Course = *CourseToResponse(&course)

		for _, class := range course.Classes {
			dataItem.Classes = append(dataItem.Classes, *ClassToResponse(&class))
		}

		dataItem.TotalStudents = (*totalStudents)[i]
		dataItems = append(dataItems, dataItem)

		// reset data class
		dataItem.Classes = nil
	}

	responses := &model.UserCourseListResponse{
		Items: dataItems,
	}

	return responses
}
