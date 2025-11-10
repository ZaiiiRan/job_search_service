package models

import (
	"slices"
)

type QueryApplicantsDal struct {
	Ids       []int64  `json:"ids"`
	Emails    []string `json:"emails"`
	IsActive  *bool    `json:"is_active"`
	IsDeleted *bool    `json:"is_deleted"`
	Limit     int      `json:"limit"`
	Offset    int      `json:"offset"`
}

func NewQueryApplicantsDal(
	ids []int64,
	emails []string,
	isActive *bool,
	isDeleted *bool,
	page, pageSize int,
) *QueryApplicantsDal {
	slices.Sort(ids)
	slices.Sort(emails)

	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}

	return &QueryApplicantsDal{
		Ids:       ids,
		Emails:    emails,
		IsActive:  isActive,
		IsDeleted: isDeleted,
		Limit:     pageSize,
		Offset:    pageSize * (page - 1),
	}
}
