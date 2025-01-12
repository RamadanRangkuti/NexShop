package dto

type CreateUserRequest struct {
	Username string `json:"name" form:"username" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username *string `json:"name" form:"username"`
	Email    *string `json:"email" form:"email"`
	Password *string `json:"password" form:"password"`
}
