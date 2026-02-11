package dto

type MarkAttendanceReq struct {
	StudentID  int    `json:"student_id" binding:"required"`
	LessonDate string `json:"lesson_date" binding:"required"`
	Status     string `json:"status" binding:"required"`
	Note       string `json:"note"`
}
