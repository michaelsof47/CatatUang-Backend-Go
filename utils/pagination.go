package utils

import (
	"math"
	"strconv"
)

type PaginationChild struct {
	Page   int
	Limit  int
	Offset int
}

type PaginationData struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalData   int `json:"total_data"`
	TotalPages  int `json:"total_pages"`
}

func ParsePagination(_page, _limit string) PaginationChild {
	page, err := strconv.Atoi(_page)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(_limit)
	if err != nil || limit < 1 {
		limit = 10
	}

	return PaginationChild{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func BuildPaginationData(params PaginationChild, totalData int) PaginationData {
	return PaginationData{
		CurrentPage: params.Page,
		Limit:       params.Limit,
		TotalData:   totalData,
		TotalPages:  int(math.Ceil(float64(totalData) / float64(params.Limit))),
	}
}
