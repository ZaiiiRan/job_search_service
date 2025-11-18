package utils

import (
	"strings"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/auth_service/v1"
)

func SanitizeRegisterApplicantRequest(req *pb.RegisterApplicantRequest) {
	req.Password = strings.TrimSpace(req.Password)
}

func SanitizeActivateApplicantRequest(req *pb.ActivateApplicantRequest) {
	req.Code = strings.TrimSpace(req.Code)
}

func SanitizeLoginApplicantRequest(req *pb.LoginApplicantRequest) {
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
}

func SanitizeGetResetApplicantPasswordCodeRequest(req *pb.GetResetApplicantPasswordCodeRequest) {
	req.Email = strings.TrimSpace(req.Email)
}

func SanitizeResetApplicantPasswordRequest(req *pb.ResetApplicantPasswordRequest) {
	req.Email = strings.TrimSpace(req.Email)
	req.Code = strings.TrimSpace(req.Code)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
}

func SanitizeChangeApplicantPasswordRequest(req *pb.ChangeApplicantPasswordRequest) {
	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
}

func SanitizeRegisterEmployerRequest(req *pb.RegisterEmployerRequest) {
	req.Password = strings.TrimSpace(req.Password)
}

func SanitizeActivateEmployerRequest(req *pb.ActivateEmployerRequest) {
	req.Code = strings.TrimSpace(req.Code)
}

func SanitizeLoginEmployerRequest(req *pb.LoginEmployerRequest) {
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
}

func SanitizeGetResetEmployerPasswordCodeRequest(req *pb.GetResetEmployerPasswordCodeRequest) {
	req.Email = strings.TrimSpace(req.Email)
}

func SanitizeResetEmployerPasswordRequest(req *pb.ResetEmployerPasswordRequest) {
	req.Email = strings.TrimSpace(req.Email)
	req.Code = strings.TrimSpace(req.Code)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
}

func SanitizeChangeEmployerPasswordRequest(req *pb.ChangeEmployerPasswordRequest) {
	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
}
