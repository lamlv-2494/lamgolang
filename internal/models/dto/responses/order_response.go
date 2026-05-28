package responses

import "time"

type OrderItemData struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	SubTotal  float64 `json:"sub_total"`
}

type OrderResponse struct {
	ID         uint            `json:"id"`
	TotalPrice float64         `json:"total_price"`
	Status     string          `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
	Items      []OrderItemData `json:"items"`
}
