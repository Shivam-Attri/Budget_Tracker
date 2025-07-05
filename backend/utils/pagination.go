// utils/pagination.go
// Provides helper functions for handling pagination logic.
package utils

import (
	"math"
	"net/http"
	"strconv"
	"your_username/budget-tracker/models"
)

// GetPaginationParams parses page and limit from query parameters.
func GetPaginationParams(r *http.Request, defaultLimit int) (page, limit, offset int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = defaultLimit
	}

	offset = (page - 1) * limit
	return page, limit, offset
}

// CalculateMetadata computes the metadata for a paginated response.
func CalculateMetadata(totalRecords, page, limit int) models.Metadata {
	if totalRecords == 0 {
		return models.Metadata{TotalRecords: 0}
	}
	return models.Metadata{
		CurrentPage:  page,
		PageSize:     limit,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(limit))),
		TotalRecords: totalRecords,
	}
}
