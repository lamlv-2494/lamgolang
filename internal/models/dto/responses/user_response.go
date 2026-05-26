package responses

type UserData struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token,omitempty"`
	Role     string `json:"role"`
}

type UserResponse struct {
	User UserData `json:"user"`
}
