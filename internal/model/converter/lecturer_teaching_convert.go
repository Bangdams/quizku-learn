package converter

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
)

func LecturerTeachingToResponse(lecturerTeaching *entity.LecturerTeaching) *model.LecturerTeachingResponse {
	log.Println("log from lecturerTeaching to response")
	return &model.LecturerTeachingResponse{
		ID:      lecturerTeaching.ID,
		Code:    lecturerTeaching.CourseCode,
		UserId:  lecturerTeaching.UserId,
		ClassId: lecturerTeaching.ClassId,
	}
}

func LecturerTeachingToResponses(lecturerTeachings *[]entity.LecturerTeaching) *[]model.LecturerTeachingResponse {
	var lecturerTeachingResponses []model.LecturerTeachingResponse

	log.Println("log from lecturerTeaching to responses")

	for _, lecturerTeaching := range *lecturerTeachings {
		lecturerTeachingResponses = append(lecturerTeachingResponses, *LecturerTeachingToResponse(&lecturerTeaching))
	}

	return &lecturerTeachingResponses
}

func DisplayDataWitClassIdToResponses(courses *[]entity.Course) *model.DisplayDataWitClassIdResponse {
	var dataItem model.DisplayDataItem
	var dataItems []model.DisplayDataItem

	log.Println("log from lecturerTeaching to response")

	for _, course := range *courses {
		if len(course.Classes) == 0 {
			continue
		}

		dataItem.Course = *CourseToResponse(&course)

		for _, class := range course.Classes {
			dataItem.Classes = append(dataItem.Classes, *ClassToResponse(&class))
		}

		dataItems = append(dataItems, dataItem)

		// reset data class
		dataItem.Classes = nil
	}

	responses := &model.DisplayDataWitClassIdResponse{
		Items: dataItems,
	}

	return responses
}

func DisplayDataToResponses(classes *[]entity.Class, courses *[]entity.Course) *model.DisplayDataResponse {
	var classResponses []model.ClassResponse
	var courseResponses []model.CourseResponse

	log.Println("log from lecturerTeaching to response")

	for _, class := range *classes {
		classResponses = append(classResponses, *ClassToResponse(&class))
	}

	for _, course := range *courses {
		courseResponses = append(courseResponses, *CourseToResponse(&course))
	}

	responses := &model.DisplayDataResponse{
		Classes: classResponses,
		Courses: courseResponses,
	}

	return responses
}
