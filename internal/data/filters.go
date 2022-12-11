package data

import "github.com/siwonpawel/greenlight/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater then zero")
	v.Check(f.Page <= 10_000, "page", "must be a maximum of 10 000")
	v.Check(f.PageSize > 0, "pageSize", "must be greater then zero")
	v.Check(f.PageSize <= 100_000, "pageSize", "must be a maximum of 100 000")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
