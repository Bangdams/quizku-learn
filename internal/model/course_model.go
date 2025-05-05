package model

type CourseResponse struct {
	CourseCode string `json:"course_code" validate:"required"`
	Name       string `json:"name" validate:"required"`
}

type CourseRequest struct {
	CourseCode string `json:"course_code" validate:"required"`
	Name       string `json:"name" validate:"required"`
}

type CourseRequestUpdate struct {
	OldCourseCode string `json:"old_course_code" validate:"required"`
	NewCourseCode string `json:"new_course_code" validate:"required"`
	Name          string `json:"name" validate:"required"`
}

type DeleteCourseRequest struct {
	CourseCode string `json:"course_code" validate:"required"`
}

type CourseRequestCourseCode struct {
	CourseCode string `json:"course_code" validate:"required"`
}
