package models

import (
	"slices"
)

type QueryApplicantsDal struct {
	Ids          []int64  `json:"ids"`
	Emails       []string `json:"emails"`
	EmailSubstrs []string `json:"email_substrs"`
	IsActive     *bool    `json:"is_active"`
	IsDeleted    *bool    `json:"is_deleted"`
	Limit        int      `json:"limit"`
	Offset       int      `json:"offset"`
}

func NewQueryApplicantsDal(
	ids []int64,
	emails []string,
	emailSubstrs []string,
	isActive *bool,
	isDeleted *bool,
	page, pageSize int,
) *QueryApplicantsDal {
	slices.Sort(ids)
	slices.Sort(emails)
	slices.Sort(emailSubstrs)

	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}

	return &QueryApplicantsDal{
		Ids:          ids,
		Emails:       emails,
		EmailSubstrs: emailSubstrs,
		IsActive:     isActive,
		IsDeleted:    isDeleted,
		Limit:        pageSize,
		Offset:       pageSize * (page - 1),
	}
}
