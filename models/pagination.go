// models/pagination.go
// Defines the structure for paginated API responses.
package models

// Metadata holds pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records"`
}

// PaginatedResponse is a generic struct for paginated API responses.
type PaginatedResponse struct {
	Metadata Metadata    `json:"metadata"`
	Data     interface{} `json:"data"`
}
