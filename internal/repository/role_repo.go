package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepo struct{ db *pgxpool.Pool }

func NewRoleRepo(db *pgxpool.Pool) *RoleRepo { return &RoleRepo{db: db} }

func (r *RoleRepo) GetNameByID(ctx context.Context, roleID int) (string, error) {
	var name string
	err := r.db.QueryRow(ctx, `SELECT name FROM roles WHERE id=$1`, roleID).Scan(&name)
	return name, err
}

func (r *RoleRepo) GetIDByName(ctx context.Context, name string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `SELECT id FROM roles WHERE name=$1`, name).Scan(&id)
	return id, err
}

func (r *RoleRepo) List(ctx context.Context) ([]struct {
	ID   int
	Name string
}, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM roles ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]struct {
		ID   int
		Name string
	}, 0)
	for rows.Next() {
		var x struct {
			ID   int
			Name string
		}
		if err := rows.Scan(&x.ID, &x.Name); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
