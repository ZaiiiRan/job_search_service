package userservice

import (
	"context"

	pb "github.com/ZaiiiRan/job_search_service/auth-service/gen/go/user_service/v1"
	usergrpcclient "github.com/ZaiiiRan/job_search_service/auth-service/internal/transport/client/grpc/user_client"
	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	CreateApplicant(ctx context.Context, applicant *pb.Applicant) (*pb.Applicant, error)
	CreateEmployer(ctx context.Context, employer *pb.Employer) (*pb.Employer, error)
	GetApplicantByEmail(ctx context.Context, email string) (*pb.Applicant, error)
	GetEmployerByEmail(ctx context.Context, email string) (*pb.Employer, error)
	GetApplicantById(ctx context.Context, id int64) (*pb.Applicant, error)
	GetEmployerById(ctx context.Context, id int64) (*pb.Employer, error)
}

type service struct {
	userClient *usergrpcclient.Client
	log        *zap.SugaredLogger
}

func New(userClient *usergrpcclient.Client, log *zap.SugaredLogger) *service {
	return &service{
		userClient: userClient,
		log:        log,
	}
}

func (s *service) CreateApplicant(ctx context.Context, applicant *pb.Applicant) (*pb.Applicant, error) {
	l := s.log.With("op", "create_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().CreateApplicant(ctx, &pb.CreateApplicantRequest{Applicant: applicant})
	if err != nil {
		l.Errorw("user.create_applicant_failed", "err", err)
		return nil, err
	}

	l.Infow("user.create_applicant.success")
	return resp.Applicant, nil
}

func (s *service) CreateEmployer(ctx context.Context, employer *pb.Employer) (*pb.Employer, error) {
	l := s.log.With("op", "create_employer", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().CreateEmployer(ctx, &pb.CreateEmployerRequest{Employer: employer})
	if err != nil {
		l.Errorw("user.create_employer_failed", "err", err)
		return nil, err
	}

	l.Infow("user.create_employer.success")
	return resp.Employer, nil
}

func (s *service) GetApplicantByEmail(ctx context.Context, email string) (*pb.Applicant, error) {
	l := s.log.With("op", "get_applicant_by_email", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().GetApplicantByEmail(ctx, &pb.GetApplicantByEmailRequest{Email: email})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			l.Warnw("user.get_applicant_by_email_failed", "err", err.Error())
			return nil, nil
		}
		l.Errorw("user.get_applicant_by_email_failed", "err", err)
		return nil, err
	}

	l.Infow("user.get_applicant_by_email.success")
	return resp.Applicant, nil
}

func (s *service) GetEmployerByEmail(ctx context.Context, email string) (*pb.Employer, error) {
	l := s.log.With("op", "get_employer_by_email", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().GetEmployerByEmail(ctx, &pb.GetEmployerByEmailRequest{Email: email})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			l.Warnw("user.get_employer_by_email_failed", "err", err.Error())
			return nil, nil
		}
		l.Errorw("user.get_employer_by_email_failed", "err", err)
		return nil, err
	}

	l.Infow("user.get_employer_by_email.success")
	return resp.Employer, nil
}

func (s *service) GetApplicantById(ctx context.Context, id int64) (*pb.Applicant, error) {
	l := s.log.With("op", "get_applicant_by_id", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().GetApplicant(ctx, &pb.GetApplicantRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			l.Warnw("user.get_applicant_by_id_failed", "err", err.Error())
			return nil, nil
		}
		l.Errorw("user.get_applicant_by_id_failed", "err", err)
		return nil, err
	}

	l.Infow("user.get_applicant_by_id.success")
	return resp.Applicant, nil
}

func (s *service) GetEmployerById(ctx context.Context, id int64) (*pb.Employer, error) {
	l := s.log.With("op", "get_employer_by_id", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().GetEmployer(ctx, &pb.GetEmployerRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			l.Warnw("user.get_employer_by_id_failed", "err", err.Error())
			return nil, nil
		}
		l.Errorw("user.get_employer_by_id_failed", "err", err)
		return nil, err
	}

	l.Infow("user.get_employer_by_id.success")
	return resp.Employer, nil
}
