package repository

import (
	"context"
	"strconv"

	"lms-backend/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AttendanceRepo struct{ db *pgxpool.Pool }

func NewAttendanceRepo(db *pgxpool.Pool) *AttendanceRepo { return &AttendanceRepo{db: db} }

func (r *AttendanceRepo) Upsert(ctx context.Context, a model.Attendance) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO attendance(course_id, student_id, lesson_date, status, note)
		 VALUES ($1,$2,$3,$4,$5)
		 ON CONFLICT (course_id, student_id, lesson_date)
		 DO UPDATE SET status = EXCLUDED.status, note = EXCLUDED.note`,
		a.CourseID, a.StudentID, a.LessonDate, a.Status, a.Note,
	)
	return err
}

func (r *AttendanceRepo) ListByCourse(ctx context.Context, courseID int) ([]model.Attendance, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, course_id, student_id, lesson_date, status, COALESCE(note,'')
		 FROM attendance
		 WHERE course_id = $1
		 ORDER BY lesson_date DESC, student_id ASC
		 LIMIT 200`,
		courseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Attendance, 0)
	for rows.Next() {
		var a model.Attendance
		if err := rows.Scan(&a.ID, &a.CourseID, &a.StudentID, &a.LessonDate, &a.Status, &a.Note); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// ListByStudent lists attendance for a student. If courseID == 0, lists across all courses.
func (r *AttendanceRepo) ListByStudent(ctx context.Context, studentID int, courseID int) ([]model.Attendance, error) {
	q := `SELECT id, course_id, student_id, lesson_date, status, COALESCE(note,'')
	      FROM attendance
	      WHERE student_id = $1`
	args := []any{studentID}
	if courseID > 0 {
		q += ` AND course_id = $2`
		args = append(args, courseID)
	}
	q += ` ORDER BY lesson_date DESC, course_id DESC LIMIT 500`

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Attendance, 0)
	for rows.Next() {
		var a model.Attendance
		if err := rows.Scan(&a.ID, &a.CourseID, &a.StudentID, &a.LessonDate, &a.Status, &a.Note); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// helper for debugging
var _ = strconv.Itoa
