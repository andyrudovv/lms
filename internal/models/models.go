package models

import "time"

type Course struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    TeacherID   *int      `json:"teacher_id"`
    SyllabusPDF string    `json:"syllabus_pdf"`
    CreatedAt   time.Time `json:"created_at"`
}

type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    RoleID    int    `json:"role_id"`
}

type Enrollment struct {
    ID         int       `json:"id"`
    UserID     int       `json:"user_id"`
    CourseID   int       `json:"course_id"`
    EnrolledAt time.Time `json:"enrolled_at"`
}

type Progress struct {
    CourseID        int     `json:"course_id"`
    StudentID       int     `json:"student_id"`
    CompletedWeeks  int     `json:"completed_weeks"`
    TotalWeeks      int     `json:"total_weeks"`
    PercentComplete float64 `json:"percent_complete"`
}