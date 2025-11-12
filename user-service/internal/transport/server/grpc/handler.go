package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/job_search_service/user-service/gen/go/user_service/v1"
	applicantservice "github.com/ZaiiiRan/job_search_service/user-service/internal/services/applicant"
	employerservice "github.com/ZaiiiRan/job_search_service/user-service/internal/services/employer"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/utils"
)

type userHandler struct {
	applicantService applicantservice.ApplicantService
	employerService  employerservice.EmployerService
	pb.UnimplementedUserServiceServer
}

func newUserHandler(applicantService applicantservice.ApplicantService, employerService employerservice.EmployerService) *userHandler {
	return &userHandler{
		applicantService: applicantService,
		employerService:  employerService,
	}
}

func (h *userHandler) CreateApplicant(ctx context.Context, req *pb.CreateApplicantRequest) (*pb.CreateApplicantResponse, error) {
	utils.SanitizeCreateApplicantRequest(req)
	return h.applicantService.CreateApplicant(ctx, req)
}

func (h *userHandler) UpdateApplicant(ctx context.Context, req *pb.UpdateApplicantRequest) (*pb.UpdateApplicantResponse, error) {
	utils.SanitizeUpdateApplicantRequest(req)
	return &pb.UpdateApplicantResponse{}, nil
}

func (h *userHandler) DeleteApplicant(ctx context.Context, req *pb.DeleteApplicantRequest) (*pb.DeleteApplicantResponse, error) {
	return &pb.DeleteApplicantResponse{}, nil
}

func (h *userHandler) QueryApplicants(ctx context.Context, req *pb.QueryApplicantsRequest) (*pb.QueryApplicantsResponse, error) {
	utils.SanitizeQueryApplicantsRequest(req)
	return h.applicantService.QueryApplicants(ctx, req)
}

func (h *userHandler) GetApplicant(ctx context.Context, req *pb.GetApplicantRequest) (*pb.GetApplicantResponse, error) {
	return h.applicantService.GetApplicant(ctx, req)
}

func (h *userHandler) GetApplicantByEmail(ctx context.Context, req *pb.GetApplicantByEmailRequest) (*pb.GetApplicantByEmailResponse, error) {
	utils.SanitizeGetApplicantByEmailRequest(req)
	return h.applicantService.GetApplicantByEmail(ctx, req)
}

func (h *userHandler) CreateEmployer(ctx context.Context, req *pb.CreateEmployerRequest) (*pb.CreateEmployerResponse, error) {
	utils.SanitizeCreateEmployerRequest(req)
	return h.employerService.CreateEmployer(ctx, req)
}

func (h *userHandler) UpdateEmployer(ctx context.Context, req *pb.UpdateEmployerRequest) (*pb.UpdateEmployerResponse, error) {
	utils.SanitizeUpdateEmployerRequest(req)
	return &pb.UpdateEmployerResponse{}, nil
}

func (h *userHandler) DeleteEmployer(ctx context.Context, req *pb.DeleteEmployerRequest) (*pb.DeleteEmployerResponse, error) {
	return &pb.DeleteEmployerResponse{}, nil
}

func (h *userHandler) QueryEmployers(ctx context.Context, req *pb.QueryEmployersRequest) (*pb.QueryEmployersResponse, error) {
	utils.SanitizeQueryEmployersRequest(req)
	return h.employerService.QueryEmployers(ctx, req)
}

func (h *userHandler) GetEmployer(ctx context.Context, req *pb.GetEmployerRequest) (*pb.GetEmployerResponse, error) {
	return h.employerService.GetEmployer(ctx, req)
}

func (h *userHandler) GetEmployerByEmail(ctx context.Context, req *pb.GetEmployerByEmailRequest) (*pb.GetEmployerByEmailResponse, error) {
	utils.SanitizeGetEmployerByEmailRequest(req)
	return h.employerService.GetEmployerByEmail(ctx, req)
}
