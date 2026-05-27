package requests

type CreateProductRequest struct {
	Name        string  `form:"name" json:"name" binding:"required"`
	Description string  `form:"description" json:"description"`
	Price       float64 `form:"price" json:"price" binding:"required,gt=0"`
	Image       string  `form:"-" json:"-"`
	CategoryID  uint    `form:"category_id" json:"category_id" binding:"required"`
}

type UpdateProductRequest struct {
	Name        *string  `form:"name" json:"name"`
	Description *string  `form:"description" json:"description"`
	Price       *float64 `form:"price" json:"price" binding:"omitempty,gt=0"`
	Image       *string  `form:"-" json:"-"`
	CategoryID  *uint    `form:"category_id" json:"category_id"`
}
