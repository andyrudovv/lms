-- +goose Up
CREATE TABLE IF NOT EXISTS roles (
  id   SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
  id            SERIAL PRIMARY KEY,
  email         TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  full_name     TEXT NOT NULL,
  role_id       INT NOT NULL REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS courses (
  id          SERIAL PRIMARY KEY,
  title       TEXT NOT NULL,
  teacher_id  INT NOT NULL REFERENCES users(id),
  created_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS enrollments (
  id          SERIAL PRIMARY KEY,
  course_id   INT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  student_id  INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  enrolled_at TIMESTAMP NOT NULL DEFAULT now(),
  UNIQUE(course_id, student_id)
);

CREATE TABLE IF NOT EXISTS attendance (
  id          SERIAL PRIMARY KEY,
  course_id   INT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  student_id  INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  lesson_date DATE NOT NULL,
  status      TEXT NOT NULL CHECK (status IN ('present','absent','late')),
  note        TEXT,
  UNIQUE(course_id, student_id, lesson_date)
);

INSERT INTO roles(name) VALUES ('admin'), ('teacher'), ('student')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS attendance;
DROP TABLE IF EXISTS enrollments;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
