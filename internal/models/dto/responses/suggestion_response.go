package responses

import "time"

type SuggestionData struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"` // Extracted from User relation
	Email       string    `json:"email"`    // Extracted from User relation
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type SuggestionListResponse struct {
	Suggestions []SuggestionData `json:"suggestions"`
}
