package models

import (
	"slices"
	"time"
)

type QueryEmployersDal struct {
	Ids                []int64    `json:"ids"`
	Emails             []string   `json:"emails"`
	EmailSubstrs       []string   `json:"email_substrs"`
	CompanyNames       []string   `json:"company_names"`
	CompanyNameSubstrs []string   `json:"company_names_substrs"`
	IsActive           *bool      `json:"is_active"`
	IsDeleted          *bool      `json:"is_deleted"`
	CreatedFrom        *time.Time `json:"created_from"`
	CreatedTo          *time.Time `json:"created_to"`
	UpdatedFrom        *time.Time `json:"updated_from"`
	UpdatedTo          *time.Time `json:"updated_to"`
	Limit              int        `json:"limit"`
	Offset             int        `json:"offset"`
}

func NewQueryEmployersDal(
	ids []int64,
	emails []string,
	emailSubstrs []string,
	companyNames []string,
	companyNameSubstrs []string,
	isActive *bool,
	isDeleted *bool,
	createdFrom, createdTo, updatedFrom, updatedTo *time.Time,
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
		CreatedFrom:        createdFrom,
		CreatedTo:          createdTo,
		UpdatedFrom:        updatedFrom,
		UpdatedTo:          updatedTo,
		Limit:              pageSize,
		Offset:             pageSize * (page - 1),
	}
}
