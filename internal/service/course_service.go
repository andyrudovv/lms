package service

import (
	"context"
	"errors"
	"strings"

	"lms-backend/internal/domain/model"
	"lms-backend/internal/repository"
)

type CourseService struct {
	courses     *repository.CourseRepo
	enrollments *repository.EnrollmentRepo
}

func NewCourseService(courses *repository.CourseRepo, enrollments *repository.EnrollmentRepo) *CourseService {
	return &CourseService{courses: courses, enrollments: enrollments}
}

func (s *CourseService) Create(ctx context.Context, c model.Course) (int, error) {
	c.Title = strings.TrimSpace(c.Title)
	if c.Title == "" {
		return 0, errors.New("title is required")
	}
	if c.TeacherID <= 0 {
		return 0, errors.New("teacher_id must be > 0")
	}
	return s.courses.Create(ctx, c)
}

func (s *CourseService) List(ctx context.Context) ([]model.Course, error) {
	return s.courses.List(ctx)
}

func (s *CourseService) ListByTeacher(ctx context.Context, teacherID int) ([]model.Course, error) {
	if teacherID <= 0 {
		return nil, errors.New("teacher_id must be > 0")
	}
	return s.courses.ListByTeacher(ctx, teacherID)
}

func (s *CourseService) ListByStudent(ctx context.Context, studentID int) ([]model.Course, error) {
	if studentID <= 0 {
		return nil, errors.New("student_id must be > 0")
	}
	return s.enrollments.ListCoursesByStudent(ctx, studentID)
}

func (s *CourseService) Enroll(ctx context.Context, courseID int, studentID int) error {
	if courseID <= 0 || studentID <= 0 {
		return errors.New("course_id and student_id must be > 0")
	}
	return s.enrollments.Enroll(ctx, courseID, studentID)
}
