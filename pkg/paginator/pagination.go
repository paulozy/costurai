package paginator

import (
	"fmt"
	"math"
)

type Paginate[T interface{}] struct {
	Items *[]T `json:"items"`
	*PaginationInfo
}

type PaginationInfo struct {
	Limit       int64  `json:"limit"`
	Page        int64  `json:"page"`
	TotalItems  int64  `json:"totalItems"`
	TotalPages  int64  `json:"totalPages"`
	NextPage    *int64 `json:"nextPage"`
	PrevPage    *int64 `json:"prevPage"`
	HasNextPage bool   `json:"hasNextPage"`
	HasPrevPage bool   `json:"hasPrevPage"`
}

func NewPaginatation(limit, page, total int64) *PaginationInfo {
	var (
		nextPage, prevPage *int64
	)

	fmt.Printf("limit %d; page %d; total pages %d", limit, page, total)

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))

	hasNextPage := page < totalPages
	hasPrevPage := page > 1

	if hasNextPage {
		value := (page + 1)
		nextPage = &value
	} else {
		nextPage = nil
	}

	if hasPrevPage {
		value := page - 1
		prevPage = &value
	} else {
		prevPage = nil
	}

	return &PaginationInfo{
		Limit:       limit,
		Page:        page,
		TotalItems:  total,
		TotalPages:  totalPages,
		NextPage:    nextPage,
		PrevPage:    prevPage,
		HasNextPage: hasNextPage,
		HasPrevPage: hasPrevPage,
	}
}

type Offset struct {
	Start int
	End   int
}

func GetOffset[T interface{}](limit, page int64, slice []T) *Offset {
	start := min(int((page-1)*limit), len(slice))
	end := min(start+int(limit), len(slice))

	return &Offset{
		Start: start,
		End:   end,
	}
}
