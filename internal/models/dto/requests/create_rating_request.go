package requests

type CreateRatingRequest struct {
	Stars   int    `json:"stars" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"max=500"`
}
