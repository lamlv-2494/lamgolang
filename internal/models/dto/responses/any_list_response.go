package responses

type AnyListResponse struct {
	Data       any   `json:"data"`
	TotalCount int64 `json:"totalCount"`
}
