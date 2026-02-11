package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"lms-backend/internal/domain/model"
	"lms-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users  *repository.UserRepo
	roles  *repository.RoleRepo
	secret []byte
	ttl    time.Duration
}

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(users *repository.UserRepo, roles *repository.RoleRepo, secret string, ttlMinutes int) *AuthService {
	return &AuthService{
		users:  users,
		roles:  roles,
		secret: []byte(secret),
		ttl:    time.Duration(ttlMinutes) * time.Minute,
	}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	if len(password) < 4 {
		return "", errors.New("password too short (min 4)")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

// RegisterStudent creates a user with role=student (no self-selected role)
func (s *AuthService) RegisterStudent(ctx context.Context, email, password, fullName string) (int, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	fullName = strings.TrimSpace(fullName)
	if email == "" || password == "" || fullName == "" {
		return 0, errors.New("email, password, full_name required")
	}

	roleID, err := s.roles.GetIDByName(ctx, "student")
	if err != nil {
		return 0, errors.New("student role not found")
	}

	hash, err := s.HashPassword(password)
	if err != nil {
		return 0, err
	}

	id, err := s.users.Create(ctx, model.User{
		Email:        email,
		PasswordHash: hash,
		FullName:     fullName,
		RoleID:       roleID,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return "", errors.New("email and password required")
	}

	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	roleName, err := s.roles.GetNameByID(ctx, u.RoleID)
	if err != nil {
		return "", errors.New("role not found")
	}

	now := time.Now()
	claims := Claims{
		UserID: u.ID,
		Role:   roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}

func (s *AuthService) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
