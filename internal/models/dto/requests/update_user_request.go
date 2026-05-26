package requests

type UpdateUserRequest struct {
	User struct {
		Username *string `json:"username,omitempty"`
		Email    *string `json:"email,omitempty"`
		Password *string `json:"password,omitempty"`
	} `json:"user"`
}
