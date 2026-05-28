package requests

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending processing delivering completed cancelled"`
}
