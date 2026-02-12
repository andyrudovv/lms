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


func (r *EnrollmentRepo) Unenroll(ctx context.Context, courseID int, studentID int) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM enrollments WHERE course_id = $1 AND student_id = $2`,
		courseID, studentID,
	)
	return err
}

func (r *EnrollmentRepo) ListEnrolledStudents(ctx context.Context, courseID int) ([]model.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT u.id, u.full_name, u.email
		 FROM enrollments e
		 JOIN users u ON u.id = e.student_id
		 WHERE e.course_id = $1
		 ORDER BY u.id DESC
		 LIMIT 200`,
		courseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.User, 0)
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *EnrollmentRepo) IsEnrolled(ctx context.Context, courseID int, studentID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM enrollments WHERE course_id = $1 AND student_id = $2)`,
		courseID, studentID,
	).Scan(&exists)
	return exists, err
}

func (r *EnrollmentRepo) ListAvailableStudents(ctx context.Context, courseID int) ([]model.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, full_name, email FROM users
		 WHERE role_id = (SELECT id FROM roles WHERE name = 'student')
		   AND id NOT IN (SELECT student_id FROM enrollments WHERE course_id = $1)
		 ORDER BY id DESC
		 LIMIT 200`,
		courseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.User, 0)
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}