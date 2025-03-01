package dto

type PaginatedResponse[T any] struct {
	Meta Metadata `json:"meta"`
	Data []T      `json:"data"`
}

type Metadata struct {
	Page        int  `json:"current_page"`
	Limit       int  `json:"limit"`
	TotalPages  int  `json:"total_pages"`
	TotalCount  int  `json:"total_count"`
	HasNextPage bool `json:"has_next_page"`
	HasPrevPage bool `json:"has_prev_page"`
}
