package dto

type CreateCourseReq struct {
	Title     string `json:"title" binding:"required"`
	TeacherID int    `json:"teacher_id"`
}

type EnrollReq struct {
	StudentID int `json:"student_id" binding:"required"`
}
