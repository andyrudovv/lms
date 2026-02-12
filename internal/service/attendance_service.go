package service

import (
	"context"
	"errors"
	"log"
	"time"

	"lms-backend/internal/domain/model"
	"lms-backend/internal/repository"
)

type AttendanceService struct {
	repo *repository.AttendanceRepo
}

func NewAttendanceService(repo *repository.AttendanceRepo) *AttendanceService {
	return &AttendanceService{repo: repo}
}

func (s *AttendanceService) Mark(ctx context.Context, a model.Attendance) error {
	switch a.Status {
	case "present", "absent", "late":
	default:
		return errors.New("status must be present|absent|late")
	}
	if a.CourseID <= 0 || a.StudentID <= 0 {
		return errors.New("course_id and student_id must be > 0")
	}
	if a.LessonDate.Equal((time.Time{})) {
		return errors.New("lesson_date is required")
	}
	return s.repo.Upsert(ctx, a)
}

func (s *AttendanceService) ListByCourse(ctx context.Context, courseID int) ([]model.Attendance, error) {
	log.Printf("Listing attendance for course ID: %d", courseID)
	if courseID <= 0 {
		return nil, errors.New("course_id must be > 0")
	}
	return s.repo.ListByCourse(ctx, courseID)
}
func (s *AttendanceService) ListByStudent(ctx context.Context, studentID int, courseID int) ([]model.Attendance, error) {
	if studentID <= 0 {
		return nil, errors.New("student_id must be > 0")
	}
	return s.repo.ListByStudent(ctx, studentID, courseID)
}
