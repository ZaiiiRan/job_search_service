package models

import "slices"

type QueryEmployersDal struct {
	Ids                []int64  `json:"ids"`
	Emails             []string `json:"emails"`
	EmailSubstrs       []string `json:"email_substrs"`
	CompanyNames       []string `json:"company_names"`
	CompanyNameSubstrs []string `json:"company_names_substrs"`
	IsActive           *bool    `json:"is_active"`
	IsDeleted          *bool    `json:"is_deleted"`
	Limit              int      `json:"limit"`
	Offset             int      `json:"offset"`
}

func NewQueryEmployersDal(
	ids []int64,
	emails []string,
	emailSubstrs []string,
	companyNames []string,
	companyNameSubstrs []string,
	isActive *bool,
	isDeleted *bool,
	page, pageSize int,
) *QueryEmployersDal {
	slices.Sort(ids)
	slices.Sort(emails)
	slices.Sort(emailSubstrs)
	slices.Sort(companyNames)
	slices.Sort(companyNameSubstrs)

	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}

	return &QueryEmployersDal{
		Ids:                ids,
		Emails:             emails,
		EmailSubstrs:       emailSubstrs,
		CompanyNames:       companyNames,
		CompanyNameSubstrs: companyNameSubstrs,
		IsActive:           isActive,
		IsDeleted:          isDeleted,
		Limit:              pageSize,
		Offset:             pageSize * (page - 1),
	}
}
