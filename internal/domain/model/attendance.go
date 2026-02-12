package model

import "time"

type Attendance struct {
	ID         int
	CourseID   int
	StudentID  int
	LessonDate time.Time // YYYY-MM-DD
	Status     string // present/absent/late
	Note       string
}
