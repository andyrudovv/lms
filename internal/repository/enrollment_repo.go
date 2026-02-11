package repository

import (
	"context"

	"lms-backend/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EnrollmentRepo struct{ db *pgxpool.Pool }

func NewEnrollmentRepo(db *pgxpool.Pool) *EnrollmentRepo { return &EnrollmentRepo{db: db} }

func (r *EnrollmentRepo) Enroll(ctx context.Context, courseID int, studentID int) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO enrollments(course_id, student_id)
		 VALUES ($1,$2) ON CONFLICT (course_id, student_id) DO NOTHING`,
		courseID, studentID,
	)
	return err
}

func (r *EnrollmentRepo) ListCoursesByStudent(ctx context.Context, studentID int) ([]model.Course, error) {
	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.title, c.teacher_id, c.created_at
		 FROM enrollments e
		 JOIN courses c ON c.id = e.course_id
		 WHERE e.student_id = $1
		 ORDER BY c.id DESC
		 LIMIT 200`,
		studentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Course, 0)
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.Title, &c.TeacherID, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
