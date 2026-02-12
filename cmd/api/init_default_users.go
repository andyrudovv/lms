package main

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)


func InitDefaultUsers(ctx context.Context, pool *pgxpool.Pool) {
	users := []struct {
		email    string
		password string
		fullName string
		role     int
	}{
		{
			email:    "admin@aitu.edu.kz",
			password: "admin123",
			fullName: "Admin User",
			role:     1,
		},
		{
			email:    "teacher@aitu.edu.kz",
			password: "teacher123",
			fullName: "Teacher User",
			role:     2,
		},
		{
			email:    "student@aitu.edu.kz",
			password: "student123",
			fullName: "Student User",
			role:     3,
		},
	}

	for _, user := range users {
		query := `INSERT INTO users(email, password_hash, full_name, role_id)
                  VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING`
        hashedPassword, err := hashPassword(user.password)
        if err != nil {
            log.Printf("Failed to hash password for user %s: %v", user.email, err)
            continue
        }
        _, err = pool.Exec(ctx, query, user.email, hashedPassword, user.fullName, user.role)

		
		if err != nil {
			log.Printf("Failed to create default user %s: %v", user.email, err)
		} else {
			log.Printf("Created default user %s with role ID %d", user.email, user.role)
		}
	}
}

func hashPassword(password string) (string, error) {
    if len(password) < 4 {
        return "", errors.New("password too short (min 4)")
    }
    b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(b), err
}
