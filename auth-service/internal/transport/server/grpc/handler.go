package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/auth_service/v1"
	authservice "github.com/ZaiiiRan/job_search_service/auth-service/internal/services/auth"
	"github.com/ZaiiiRan/job_search_service/auth-service/internal/utils"
)

type authHandler struct {
	pb.UnimplementedAuthServiceServer
	authService authservice.AuthService
}

func newAuthHandler(authService authservice.AuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

func (h *authHandler) RegisterApplicant(ctx context.Context, req *pb.RegisterApplicantRequest) (*pb.RegisterApplicantResponse, error) {
	utils.SanitizeRegisterApplicantRequest(req)
	return h.authService.RegisterApplicant(ctx, req)
}

func (h *authHandler) GetNewApplicantActivationCode(ctx context.Context, req *pb.GetNewApplicantActivationCodeRequest) (*pb.GetNewApplicantActivationCodeResponse, error) {
	return nil, nil
}

func (h *authHandler) ActivateApplicant(ctx context.Context, req *pb.ActivateApplicantRequest) (*pb.ActivateApplicantResponse, error) {
	utils.SanitizeActivateApplicantRequest(req)
	return nil, nil
}

func (h *authHandler) LoginApplicant(ctx context.Context, req *pb.LoginApplicantRequest) (*pb.LoginApplicantResponse, error) {
	utils.SanitizeLoginApplicantRequest(req)
	return nil, nil
}

func (h *authHandler) RefreshApplicant(ctx context.Context, req *pb.RefreshApplicantRequest) (*pb.RefreshApplicantResponse, error) {
	return nil, nil
}

func (h *authHandler) LogoutApplicant(ctx context.Context, req *pb.LogoutApplicantRequest) (*pb.LogoutApplicantResponse, error) {
	return nil, nil
}

func (h *authHandler) GetResetApplicantPasswordCode(ctx context.Context, req *pb.GetResetApplicantPasswordCodeRequest) (*pb.GetResetApplicantPasswordCodeResponse, error) {
	utils.SanitizeGetResetApplicantPasswordCodeRequest(req)
	return nil, nil
}

func (h *authHandler) ResetApplicantPassword(ctx context.Context, req *pb.ResetApplicantPasswordRequest) (*pb.ResetApplicantPasswordResponse, error) {
	utils.SanitizeResetApplicantPasswordRequest(req)
	return nil, nil
}

func (h *authHandler) ChangeApplicantPassword(ctx context.Context, req *pb.ChangeApplicantPasswordRequest) (*pb.ChangeApplicantPasswordResponse, error) {
	utils.SanitizeChangeApplicantPasswordRequest(req)
	return nil, nil
}

func (h *authHandler) RegisterEmployer(ctx context.Context, req *pb.RegisterEmployerRequest) (*pb.RegisterEmployerResponse, error) {
	utils.SanitizeRegisterEmployerRequest(req)
	return nil, nil
}

func (h *authHandler) GetNewEmployerActivationCode(ctx context.Context, req *pb.GetNewEmployerActivationCodeRequest) (*pb.GetNewEmployerActivationCodeResponse, error) {
	return nil, nil
}

func (h *authHandler) ActivateEmployer(ctx context.Context, req *pb.ActivateEmployerRequest) (*pb.ActivateEmployerResponse, error) {
	utils.SanitizeActivateEmployerRequest(req)
	return nil, nil
}

func (h *authHandler) LoginEmployer(ctx context.Context, req *pb.LoginEmployerRequest) (*pb.LoginEmployerResponse, error) {
	utils.SanitizeLoginEmployerRequest(req)
	return nil, nil
}

func (h *authHandler) RefreshEmployer(ctx context.Context, req *pb.RefreshEmployerRequest) (*pb.RefreshEmployerResponse, error) {
	return nil, nil
}

func (h *authHandler) LogoutEmployer(ctx context.Context, req *pb.LogoutEmployerRequest) (*pb.LogoutEmployerResponse, error) {
	return nil, nil
}

func (h *authHandler) GetResetEmployerPasswordCode(ctx context.Context, req *pb.GetResetEmployerPasswordCodeRequest) (*pb.GetResetEmployerPasswordCodeResponse, error) {
	utils.SanitizeGetResetEmployerPasswordCodeRequest(req)
	return nil, nil
}

func (h *authHandler) ResetEmployerPassword(ctx context.Context, req *pb.ResetEmployerPasswordRequest) (*pb.ResetEmployerPasswordResponse, error) {
	utils.SanitizeResetEmployerPasswordRequest(req)
	return nil, nil
}

func (h *authHandler) ChangeEmployerPassword(ctx context.Context, req *pb.ChangeEmployerPasswordRequest) (*pb.ChangeEmployerPasswordResponse, error) {
	utils.SanitizeChangeEmployerPasswordRequest(req)
	return nil, nil
}
