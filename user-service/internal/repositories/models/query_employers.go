package models

import "slices"

type QueryEmployersDal struct {
	Ids          []int64  `json:"ids"`
	Emails       []string `json:"emails"`
	CompanyNames []string `json:"company_names"`
	IsActive     *bool    `json:"is_active"`
	IsDeleted    *bool    `json:"is_deleted"`
	Limit        int      `json:"limit"`
	Offset       int      `json:"offset"`
}

func NewQueryEmployersDal(
	ids []int64,
	emails []string,
	companyNames []string,
	isActive *bool,
	isDeleted *bool,
	page, pageSize int,
) *QueryEmployersDal {
	slices.Sort(ids)
	slices.Sort(emails)
	slices.Sort(companyNames)

	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}

	return &QueryEmployersDal{
		Ids:          ids,
		Emails:       emails,
		CompanyNames: companyNames,
		IsActive:     isActive,
		IsDeleted:    isDeleted,
		Limit:        pageSize,
		Offset:       pageSize * (page - 1),
	}
}
