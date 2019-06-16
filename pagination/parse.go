package pagination

import "math"

const (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Parse adjust pageSize and currentPage to respect `MaxPageSize`.
//
// `currentPage` starts at 1.
func Parse(pageSizeRaw, currentPageRaw int) (pageSize int, currentPage int) {
	pageSize = int(math.Min(float64(pageSizeRaw), float64(MaxPageSize)))
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}

	currentPage = int(math.Max(1, float64(currentPageRaw)))

	return pageSize, currentPage
}
