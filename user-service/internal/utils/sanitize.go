package utils

import (
	"strings"

	pb "github.com/ZaiiiRan/job_search_service/user-service/gen/go/user_service/v1"
)

func SanitizeContacts(contacts *pb.Contacts) {
	if contacts != nil {
		if contacts.PhoneNumber != nil {
			phoneNumber := strings.TrimSpace(*contacts.PhoneNumber)
			contacts.PhoneNumber = &phoneNumber
		}
		if contacts.Telegram != nil {
			telegram := strings.TrimSpace(*contacts.Telegram)
			contacts.Telegram = &telegram
		}
	}
}

func SanitizeApplicant(applicant *pb.Applicant) {
	if applicant != nil {
		applicant.FirstName = strings.TrimSpace(applicant.FirstName)
		applicant.LastName = strings.TrimSpace(applicant.LastName)
		if applicant.Patronymic != nil {
			patronymic := strings.TrimSpace(*applicant.Patronymic)
			applicant.Patronymic = &patronymic
		}
		applicant.BirthDate = strings.TrimSpace(applicant.BirthDate)
		applicant.City = strings.TrimSpace(applicant.City)
		applicant.Email = strings.TrimSpace(applicant.Email)
		SanitizeContacts(applicant.Contacts)
	}
}

func SanitizeCreateApplicantRequest(req *pb.CreateApplicantRequest) {
	SanitizeApplicant(req.Applicant)
}

func SanitizeUpdateApplicantRequest(req *pb.UpdateApplicantRequest) {
	SanitizeApplicant(req.Applicant)
}

func SanitizeQueryApplicantsRequest(req *pb.QueryApplicantsRequest) {
	for i, email := range req.FullEmails {
		req.FullEmails[i] = strings.TrimSpace(email)
	}
	for i, emailSubstr := range req.SubstrEmails {
		req.SubstrEmails[i] = strings.TrimSpace(emailSubstr)
	}
}

func SanitizeGetApplicantByEmailRequest(req *pb.GetApplicantByEmailRequest) {
	req.Email = strings.TrimSpace(req.Email)
}

func SanitizeEmployer(employer *pb.Employer) {
	if employer != nil {
		employer.CompanyName = strings.TrimSpace(employer.CompanyName)
		employer.City = strings.TrimSpace(employer.City)
		employer.Email = strings.TrimSpace(employer.Email)
		SanitizeContacts(employer.Contacts)
	}
}

func SanitizeCreateEmployerRequest(req *pb.CreateEmployerRequest) {
	SanitizeEmployer(req.Employer)
}

func SanitizeUpdateEmployerRequest(req *pb.UpdateEmployerRequest) {
	SanitizeEmployer(req.Employer)
}

func SanitizeQueryEmployersRequest(req *pb.QueryEmployersRequest) {
	for i, email := range req.FullEmails {
		req.FullEmails[i] = strings.TrimSpace(email)
	}
	for i, companyName := range req.FullCompanyNames {
		req.FullCompanyNames[i] = strings.TrimSpace(companyName)
	}
	for i, emailSubstr := range req.SubstrEmails {
		req.SubstrEmails[i] = strings.TrimSpace(emailSubstr)
	}
	for i, companyNameSubstr := range req.SubstrCompanyNames {
		req.SubstrCompanyNames[i] = strings.TrimSpace(companyNameSubstr)
	}
}

func SanitizeGetEmployerByEmailRequest(req *pb.GetEmployerByEmailRequest) {
	req.Email = strings.TrimSpace(req.Email)
}
