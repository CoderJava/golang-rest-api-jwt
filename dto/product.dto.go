package dto

type CreateProductRequest struct {
	Name string `json:"name" form:"name" binding:"required,min=1"`
	Price uint64 `json:"price" form:"price" binding:"required"`
}