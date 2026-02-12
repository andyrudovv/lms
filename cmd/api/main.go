package main

import (
	"context"
	"fmt"
	"log"

	"lms-backend/internal/config"
	"lms-backend/internal/db"
	"lms-backend/internal/repository"
	"lms-backend/internal/service"
	httpapi "lms-backend/internal/transport/http"
	"lms-backend/internal/transport/http/handlers"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("config load error: ", err)
	}

	if cfg.Migrations.AutoUp {
		if err := db.RunMigrations(cfg.DB.DSN, cfg.Migrations.Dir); err != nil {
			log.Fatal("migration error: ", err)
		}
		log.Println("migrations: up to date")
	}

	pool, err := db.NewPostgresPool(db.Options{
		DSN:                cfg.DB.DSN,
		MaxConns:           cfg.DB.MaxConns,
		MinConns:           cfg.DB.MinConns,
		MaxConnIdleMinutes: cfg.DB.MaxConnIdleMinutes,
	})
	if err != nil {
		log.Fatal("db connect error: ", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepo(pool)
	roleRepo := repository.NewRoleRepo(pool)
	courseRepo := repository.NewCourseRepo(pool)
	enrollRepo := repository.NewEnrollmentRepo(pool)
	attRepo := repository.NewAttendanceRepo(pool)

	authSvc := service.NewAuthService(userRepo, roleRepo, cfg.JWT.Secret, cfg.JWT.AccessTTLMinutes)
	userSvc := service.NewUserService(userRepo, roleRepo, authSvc)
	courseSvc := service.NewCourseService(courseRepo, enrollRepo)
	attSvc := service.NewAttendanceService(attRepo)

	authH := handlers.NewAuthHandler(authSvc)
	userH := handlers.NewUserHandler(userSvc)
	courseH := handlers.NewCourseHandler(courseSvc)
	attH := handlers.NewAttendanceHandler(attSvc)

	InitDB(context.Background(), pool) // Initialize database tables and default roles
	InitDefaultUsers(context.Background(), pool) // Initialize default users before starting the server
	
	r := httpapi.NewRouter(authSvc, authH, userH, courseH, attH)

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Println("API listening on", addr)
	log.Fatal(r.Run(addr))
}
