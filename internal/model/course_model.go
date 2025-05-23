package model

type CourseResponse struct {
	CourseCode string `json:"course_code" validate:"required"`
	Name       string `json:"name" validate:"required"`
}

type UserCourseListItem struct {
	Course        CourseResponse  `json:"course" validate:"required"`
	Classes       []ClassResponse `json:"classes" validate:"required"`
	TotalStudents uint            `json:"total_students" validate:"required"`
}

type UserCourseListResponse struct {
	Items []UserCourseListItem `json:"items" validate:"required"`
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

type CourseRequestCourseCode struct {
	CourseCode string `json:"course_code" validate:"required"`
}
