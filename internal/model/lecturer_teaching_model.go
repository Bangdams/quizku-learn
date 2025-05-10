package model

type LecturerTeachingResponse struct {
	ID      uint   `json:"id" validate:"required"`
	Code    string `json:"code" validate:"required"`
	UserId  uint   `json:"user_id" validate:"required"`
	ClassId uint   `json:"class_id" validate:"required"`
}

type DisplayDataResponse struct {
	Classes []ClassResponse  `json:"classes" validate:"required"`
	Courses []CourseResponse `json:"courses" validate:"required"`
}

type DisplayDataItem struct {
	Course  CourseResponse  `json:"course" validate:"required"`
	Classes []ClassResponse `json:"classes" validate:"required"`
}

type DisplayDataWitClassIdResponse struct {
	Items []DisplayDataItem `json:"items" validate:"required"`
}

type TeachingItem struct {
	ClassID     uint     `json:"class_id" validate:"required"`
	CourseCodes []string `json:"course_codes" validate:"required"`
}

type FlexibleLecturerTeachingRequest struct {
	UserID    uint           `json:"user_id" validate:"required"`
	Teachings []TeachingItem `json:"teachings" validate:"required"`
}
