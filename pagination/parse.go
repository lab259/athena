package pagination

import "math"

var (
	defaultPageSize = 10
	maxPageSize     = 100
)

type Pagination struct {
	CurrentPage int
	PageSize    int
	Offset      int
	Limit       int
}

// Parse adjust pageSize and currentPage to respect default page size
// and max page size. Keep in mind that the first page is `1`, therefore non-positive
// numbers will always return `1`.
func Parse(pageSizeRaw, currentPageRaw int) *Pagination {
	pageSize := int(math.Min(float64(pageSizeRaw), float64(maxPageSize)))
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	currentPage := int(math.Max(1, float64(currentPageRaw)))
	return &Pagination{
		PageSize:    pageSize,
		Limit:       pageSize,
		CurrentPage: currentPage,
		Offset:      (currentPage - 1) * pageSize,
	}
}

// Option represents functions that can be used to customize
// pagination defaults.
//
//
// See Also
//
// pagination.Setup
type Option func()

// MaxPageSize customize max page size.
func MaxPageSize(n int) Option {
	return func() {
		maxPageSize = n
	}
}

// DefaultPageSize customize default page size.
func DefaultPageSize(n int) Option {
	return func() {
		defaultPageSize = n
	}
}

// Setup allows customize default values:
//  pagination.Setup(pagination.DefaultPageSize(20), pagination.MaxPageSize(50))
func Setup(options ...Option) {
	for _, opt := range options {
		opt()
	}
}
