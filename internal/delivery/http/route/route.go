package route

import (
	"github.com/Bangdams/quizku-learn/internal/delivery/http"
	"github.com/Bangdams/quizku-learn/internal/util"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                        *fiber.App
	UserController             http.UserController
	CourseController           http.CourseController
	ClassController            http.ClassController
	LecturerTeachingController http.LecturerTeachingController
	QuestionController         http.QuestionController
	QuizController             http.QuizController
	UserAnswerController       http.UserAnswerController
}

func (config *RouteConfig) Setup() {
	// API ADMIN
	admin := config.App.Group("/api-admin", util.CheckLevel("admin"))

	// API for user
	admin.Get("/users", config.UserController.FindAll)
	admin.Get("/users/:email", config.UserController.FindByEmail)
	admin.Post("/users", config.UserController.Create)
	admin.Delete("/users/:id", config.UserController.Delete)
	admin.Put("/users", config.UserController.Update)

	// show dashboard
	admin.Get("/admin-dashboard", config.UserController.AdminDashboardReport)

	// API for course
	admin.Get("/courses", config.CourseController.FindAll)
	admin.Get("/courses/:course_code", config.CourseController.FindByCourseCode)
	admin.Post("/courses", config.CourseController.Create)
	admin.Delete("/courses/:course_code", config.CourseController.Delete)
	admin.Put("/courses", config.CourseController.Update)

	// API for question
	admin.Get("/questions", config.QuestionController.FindAll)
	admin.Delete("/questions/:question_id", config.QuestionController.Delete)

	// API for class
	admin.Get("/classes", config.ClassController.FindAll)
	admin.Get("/classes/:class_name", config.ClassController.FindByName)
	admin.Post("/classes", config.ClassController.Create)
	admin.Delete("/classes/:class_id", config.ClassController.Delete)
	admin.Put("/classes", config.ClassController.Update)

	// API for class subject
	admin.Post("/class-subject", config.ClassController.ClassSubject)

	// API for Lecturer Teaching
	admin.Get("/lecturer-teachings/:lecturer_teaching_id", config.LecturerTeachingController.FindById)
	admin.Post("/lecturer-teachings", config.LecturerTeachingController.Create)
	admin.Delete("/lecturer-teachings/:lecturer_teaching_id", config.LecturerTeachingController.Delete)

	// display Lecturer Teaching for insert
	admin.Get("/lecturer-teachings", config.LecturerTeachingController.DisplayData)

	// API DOSEN
	lecturer := config.App.Group("/api-lecturer", util.CheckLevel("dosen"))

	// API question
	lecturer.Post("/questions", config.QuestionController.Create)
	lecturer.Get("/question/:course_code", config.QuestionController.FindByCourseCode)

	// API courses
	lecturer.Get("/courses-by-user-with-class", config.CourseController.ListCoursesByUserWithClass)

	// API show dashboard
	lecturer.Get("/lecturer-dashboard", config.UserController.LecturerDashboardReport)

	// API quiz
	lecturer.Get("/quizz-dashboard", config.QuizController.QuizDashboard)
	lecturer.Post("/quizz", config.QuizController.Create)
	lecturer.Delete("/quizz/:quiz_id", config.QuizController.Delete)

	// API MAHASISWA
	student := config.App.Group("/api-student", util.CheckLevel("mahasiswa"))

	// API course
	student.Get("/courses-by-user", config.CourseController.ListCoursesByUser)

	// API quiz
	student.Get("/quizzes/course/:course_code", config.QuizController.FindByUserAndCourse)
	student.Get("/quizzes/:quiz_id/student/result", config.QuizController.QuizStudentResult)

	// API user answer
	student.Post("/quizzes/answer/:quiz_id", config.UserAnswerController.Create)

	// Api for login
	config.App.Post("/login", config.UserController.Login)
	config.App.Post("/logout", config.UserController.Logout)
	config.App.Post("/refresh", config.UserController.Refresh)
	config.App.Get("/api/status-login", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"message": "success"})
	})
}
