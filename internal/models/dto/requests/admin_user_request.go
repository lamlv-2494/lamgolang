package requests

type AdminUpdateUserRequest struct {
	User struct {
		Username *string `json:"username,omitempty"`
		Email    *string `json:"email,omitempty"`
		Role     *string `json:"role,omitempty"`
	} `json:"user"`
}
