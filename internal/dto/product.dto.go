package dto

type CreateProductRequest struct {
	Name        string  `json:"name" form:"name" binding:"required"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price" binding:"required" `
	Stock       int     `json:"stock" form:"stock" binding:"required"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name" form:"name"`
	Description *string  `json:"description" form:"description"`
	Price       *float64 `json:"price" form:"price" `
	Stock       *int     `json:"stock" form:"stock"`
}
