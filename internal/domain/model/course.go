package model

import "time"

type Course struct {
	ID        int
	Title     string
	TeacherID int
	CreatedAt time.Time
}
