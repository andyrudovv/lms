package model

type Attendance struct {
	ID         int
	CourseID   int
	StudentID  int
	LessonDate string // YYYY-MM-DD
	Status     string // present/absent/late
	Note       string
}
