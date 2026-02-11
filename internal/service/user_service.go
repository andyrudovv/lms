package service

import (
	"context"
	"errors"
	"strings"

	"lms-backend/internal/domain/model"
	"lms-backend/internal/repository"
)

type UserService struct {
	repo  *repository.UserRepo
	roles *repository.RoleRepo
	auth  *AuthService
}

func NewUserService(repo *repository.UserRepo, roles *repository.RoleRepo, auth *AuthService) *UserService {
	return &UserService{repo: repo, roles: roles, auth: auth}
}

func (s *UserService) Create(ctx context.Context, u model.User, rawPassword string) (int, error) {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.FullName = strings.TrimSpace(u.FullName)

	if u.Email == "" || u.FullName == "" {
		return 0, errors.New("email and full_name are required")
	}
	if u.RoleID <= 0 {
		return 0, errors.New("role_id must be > 0")
	}

	hash, err := s.auth.HashPassword(rawPassword)
	if err != nil {
		return 0, err
	}
	u.PasswordHash = hash

	return s.repo.Create(ctx, u)
}

func (s *UserService) List(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) Me(ctx context.Context, userID int) (model.User, string, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return model.User{}, "", err
	}
	roleName, err := s.roles.GetNameByID(ctx, u.RoleID)
	if err != nil {
		return model.User{}, "", err
	}
	return u, roleName, nil
}

// Admin only: change role by role name
func (s *UserService) ChangeRole(ctx context.Context, userID int, roleName string) error {
	roleName = strings.TrimSpace(strings.ToLower(roleName))
	if roleName == "" {
		return errors.New("role is required")
	}
	switch roleName {
	case "admin", "teacher", "student":
	default:
		return errors.New("role must be admin|teacher|student")
	}
	roleID, err := s.roles.GetIDByName(ctx, roleName)
	if err != nil {
		return err
	}
	return s.repo.UpdateRole(ctx, userID, roleID)
}

func (s *UserService) ListRoles(ctx context.Context) ([]struct{ ID int; Name string }, error) {
	return s.roles.List(ctx)
}
