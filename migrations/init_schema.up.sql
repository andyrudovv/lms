-- =========================
-- ROLES
-- =========================
CREATE TABLE roles
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- =========================
-- USERS
-- =========================
CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id       INTEGER NOT NULL REFERENCES roles,
    created_at    TIMESTAMP DEFAULT now(),
    updated_at    TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_role_id ON users (role_id);

-- =========================
-- COURSES
-- =========================
CREATE TABLE courses
(
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    teacher_id   INTEGER REFERENCES users,
    syllabus_pdf VARCHAR(255),
    created_at   TIMESTAMP DEFAULT now(),
    updated_at   TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_courses_teacher_id ON courses (teacher_id);

-- =========================
-- WEEKS
-- =========================
CREATE TABLE weeks
(
    id          SERIAL PRIMARY KEY,
    course_id   INTEGER REFERENCES courses,
    week_number INTEGER NOT NULL,
    title       VARCHAR(255)
);

CREATE INDEX idx_weeks_course_id ON weeks (course_id);

-- =========================
-- ENROLLMENTS
-- =========================
CREATE TABLE enrollments
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER REFERENCES users,
    course_id   INTEGER REFERENCES courses,
    enrolled_at TIMESTAMP DEFAULT now(),
    UNIQUE (user_id, course_id)
);

CREATE UNIQUE INDEX idx_enrollments_user_course
    ON enrollments (user_id, course_id);

CREATE INDEX idx_enrollments_course_id ON enrollments (course_id);
CREATE INDEX idx_enrollments_user_id ON enrollments (user_id);

-- =========================
-- ATTENDANCE
-- =========================
CREATE TABLE attendance
(
    id         SERIAL PRIMARY KEY,
    course_id  INTEGER REFERENCES courses,
    student_id INTEGER REFERENCES users,
    date       DATE NOT NULL,
    status     VARCHAR(20) DEFAULT 'present' NOT NULL,
    notes      TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (course_id, student_id, date)
);

CREATE UNIQUE INDEX idx_attendance_course_student_date
    ON attendance (course_id, student_id, date);

CREATE INDEX idx_attendance_course_id ON attendance (course_id);
CREATE INDEX idx_attendance_student_id ON attendance (student_id);
CREATE INDEX idx_attendance_date ON attendance (date);

-- =========================
-- EXAMS
-- =========================
CREATE TABLE exams
(
    id          SERIAL PRIMARY KEY,
    course_id   INTEGER REFERENCES courses,
    title       VARCHAR(255) NOT NULL,
    type        VARCHAR(50) NOT NULL,
    description TEXT,
    due_date    TIMESTAMP,
    locked      BOOLEAN DEFAULT false,
    created_at  TIMESTAMP DEFAULT now(),
    updated_at  TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_exams_course_id ON exams (course_id);
CREATE INDEX idx_exams_type ON exams (type);

-- =========================
-- SUBMISSIONS
-- =========================
CREATE TABLE submissions
(
    id           SERIAL PRIMARY KEY,
    exam_id      INTEGER REFERENCES exams,
    student_id   INTEGER REFERENCES users,
    file_path    VARCHAR(255),
    submitted_at TIMESTAMP DEFAULT now(),
    grade        INTEGER,
    feedback     TEXT,
    UNIQUE (exam_id, student_id)
);

CREATE UNIQUE INDEX idx_submissions_exam_student
    ON submissions (exam_id, student_id);

CREATE INDEX idx_submissions_exam_id ON submissions (exam_id);
CREATE INDEX idx_submissions_student_id ON submissions (student_id);

-- =========================
-- RESOURCES
-- =========================
CREATE TABLE resources
(
    id          SERIAL PRIMARY KEY,
    week_id     INTEGER REFERENCES weeks,
    title       VARCHAR(255) NOT NULL,
    type        VARCHAR(50) NOT NULL,
    url         VARCHAR(500),
    description TEXT,
    created_at  TIMESTAMP DEFAULT now(),
    updated_at  TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_resources_week_id ON resources (week_id);
CREATE INDEX idx_resources_type ON resources (type);

-- =========================
-- ANNOUNCEMENTS
-- =========================
CREATE TABLE announcements
(
    id         SERIAL PRIMARY KEY,
    course_id  INTEGER REFERENCES courses,
    title      VARCHAR(255),
    content    TEXT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_announcements_course_id
    ON announcements (course_id);
