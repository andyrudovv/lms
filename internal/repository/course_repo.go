package repository

import (
	"context"

	"lms-backend/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CourseRepo struct{ db *pgxpool.Pool }

func NewCourseRepo(db *pgxpool.Pool) *CourseRepo { return &CourseRepo{db: db} }

func (r *CourseRepo) Create(ctx context.Context, c model.Course) (int, error) {
	var id int
	err := r.db.QueryRow(ctx,
		`INSERT INTO courses(title, teacher_id) VALUES ($1,$2) RETURNING id`,
		c.Title, c.TeacherID,
	).Scan(&id)
	return id, err
}

func (r *CourseRepo) List(ctx context.Context) ([]model.Course, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, teacher_id, created_at FROM courses ORDER BY id DESC LIMIT 200`)
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

func (r *CourseRepo) ListByTeacher(ctx context.Context, teacherID int) ([]model.Course, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, teacher_id, created_at FROM courses WHERE teacher_id=$1 ORDER BY id DESC LIMIT 200`, teacherID)
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
