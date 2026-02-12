# LMS Backend (Go/Gin/Postgres)

## Run
1) Ensure Postgres is running and DB exists: `aitu_lms`
2) Edit `config.yaml` if needed.
3) From project root:
```bash
go mod tidy
go run ./cmd/api
```

## Auth
- POST /api/v1/auth/register  -> creates user with role=student
- POST /api/v1/auth/login     -> returns JWT

JWT header:
`Authorization: Bearer <token>`

## Profile
- GET /api/v1/me              -> current user profile from token

## Admin
- GET /api/v1/roles           -> list roles
- GET /api/v1/users           -> list users
- POST /api/v1/users          -> create user with role_id
- PATCH /api/v1/users/:id/role -> change role by name {"role":"teacher"}

## Courses
- GET /api/v1/courses         -> all courses
- POST /api/v1/courses        -> create (admin/teacher)
- GET /api/v1/my/courses      -> my courses by role
- POST /api/v1/courses/:id/enroll -> enroll student (admin/teacher)

## Attendance
- POST /api/v1/courses/:id/attendance -> mark attendance (admin/teacher)
- GET /api/v1/courses/:id/attendance  -> list course attendance (admin/teacher)
- GET /api/v1/my/attendance?course_id= -> student attendance (by token)

  
---

# Frontend (ReactJS)

```
cd frontend
npm i
npm run dev
```