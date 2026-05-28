package responses

type RatingResponseData struct {
	ProductID uint   `json:"product_id,omitempty"`
	Stars     int    `json:"stars,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

type RatingResponse struct {
	Rating RatingResponseData `json:"rating"`
}

type RatingListResponse struct {
	Ratings []RatingResponseData `json:"ratings"`
}
