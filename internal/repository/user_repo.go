package repository

import (
	"context"
	"errors"

	"lms-backend/internal/domain/model"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u model.User) (int, error) {
	var id int
	err := r.db.QueryRow(ctx,
		`INSERT INTO users(email, password_hash, full_name, role_id)
		 VALUES ($1,$2,$3,$4) RETURNING id`,
		u.Email, u.PasswordHash, u.FullName, u.RoleID,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, errors.New("email already registered")
		}
		return 0, err
	}
	return id, nil
}

func (r *UserRepo) List(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Query(ctx, `SELECT id, email, password_hash, full_name, role_id FROM users ORDER BY id DESC LIMIT 200`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.User, 0)
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, full_name, role_id FROM users WHERE email=$1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID)
	return u, err
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, full_name, role_id FROM users WHERE id=$1`,
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID)
	return u, err
}

func (r *UserRepo) UpdateRole(ctx context.Context, userID int, roleID int) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET role_id=$1 WHERE id=$2`, roleID, userID)
	return err
}