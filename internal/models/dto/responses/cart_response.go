package responses

type CartItemData struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Image     string  `json:"image"`
	Quantity  int     `json:"quantity"`
	SubTotal  float64 `json:"sub_total"`
}

type CartResponse struct {
	Items       []CartItemData `json:"items"`
	TotalAmount float64        `json:"total_amount"`
}
