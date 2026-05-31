package responses

// Struct gọn nhẹ cho Category, loại bỏ hoàn toàn các trường thời gian thừa
type CategoryCompact struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Cấu trúc sản phẩm sạch sẽ trả về cho Client
type ProductData struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Image       string          `json:"image"`
	Category    CategoryCompact `json:"category"`
	Rating      float64         `json:"rating"`
}
