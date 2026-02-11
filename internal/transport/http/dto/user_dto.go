package dto

type CreateUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	RoleID   int    `json:"role_id" binding:"required"`
}

type ChangeRoleReq struct {
	Role string `json:"role" binding:"required"`
}
