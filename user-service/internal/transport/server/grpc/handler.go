package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/job_search_service/user-service/gen/go/user_service/v1"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/utils"
)

type userHandler struct {
	pb.UnimplementedUserServiceServer
}

func newUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) ApplicantBatchCreate(ctx context.Context, req *pb.ApplicantBatchCreateRequest) (*pb.ApplicantBatchCreateResponse, error) {
	utils.SanitizeApplicantBatchCreateRequest(req)
	return &pb.ApplicantBatchCreateResponse{}, nil
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
	return &pb.QueryApplicantsResponse{}, nil
}

func (h *userHandler) EmployerBatchCreate(ctx context.Context, req *pb.EmployerBatchCreateRequest) (*pb.EmployerBatchCreateResponse, error) {
	utils.SanitizeEmployerBatchCreateRequest(req)
	return &pb.EmployerBatchCreateResponse{}, nil
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
	return &pb.QueryEmployersResponse{}, nil
}
