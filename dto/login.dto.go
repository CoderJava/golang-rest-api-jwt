package dto

type LoginRequest struct {
	Email   string `json:"email" form:"email" binding:"required,email"`
	Pasword string `json:"password" form:"password" binding:"required,min=6"`
}
