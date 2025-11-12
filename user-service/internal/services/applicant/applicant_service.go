package applicantservice

import (
	"context"
	"fmt"

	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	pb "github.com/ZaiiiRan/job_search_service/user-service/gen/go/user_service/v1"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/applicant"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/errors/validationerror"
	dal "github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/transport/postgres"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/transport/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ApplicantService interface {
	CreateApplicant(ctx context.Context, req *pb.CreateApplicantRequest) (*pb.CreateApplicantResponse, error)
	GetApplicant(ctx context.Context, req *pb.GetApplicantRequest) (*pb.GetApplicantResponse, error)
	GetApplicantByEmail(ctx context.Context, req *pb.GetApplicantByEmailRequest) (*pb.GetApplicantByEmailResponse, error)
	QueryApplicants(ctx context.Context, req *pb.QueryApplicantsRequest) (*pb.QueryApplicantsResponse, error)
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *applicantDataProvider
}

func New(pgClient *postgres.PostgresClient, redisClient *redis.RedisClient, log *zap.SugaredLogger) ApplicantService {
	return &service{
		dataProvider: newApplicantDataProvider(pgClient, redisClient),
		log:          log,
	}
}

func (s *service) CreateApplicant(ctx context.Context, req *pb.CreateApplicantRequest) (*pb.CreateApplicantResponse, error) {
	l := s.log.With("op", "create_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	a, verr := s.createApplicant(req)
	if len(verr) > 0 {
		l.Errorw("applicant.create_applicant_failed.validation_error", "err", verr)
		return nil, verr.ToStatus()
	}

	existed, err := s.dataProvider.GetByEmail(ctx, a.Email())
	if err != nil {
		l.Errorw("applicant.create_applicant_failed.check_existing_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if existed != nil {
		if existed.IsActive() {
			l.Errorw("applicant.create_applicant_failed", "err", err)
			return nil, status.Errorf(codes.AlreadyExists, "applicant with this email already exists")
		}
		l.Infow("applicant.create_applicant.restoring_inactive_applicant", "id", existed.Id())
		a.SetId(existed.Id())
	}

	if err := s.dataProvider.Save(ctx, a); err != nil {
		l.Errorw("applicant.create_applicant_failed.save_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	l.Infow("applicant.create_applicant.created")
	return &pb.CreateApplicantResponse{Applicant: toPbApplicant(a)}, nil
}

func (s *service) GetApplicant(ctx context.Context, req *pb.GetApplicantRequest) (*pb.GetApplicantResponse, error) {
	l := s.log.With("op", "get_applicant", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "id", req.Id)

	if req.Id < 1 {
		l.Errorw("applicant.get_applicant_failed", "err", "id must be positive")
		return nil, status.Errorf(codes.InvalidArgument, "id must be positive")
	}

	a, err := s.dataProvider.GetById(ctx, req.Id)
	if err != nil {
		l.Errorw("applicant.get_applicant_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if a == nil {
		return nil, status.Errorf(codes.NotFound, "applicant not found")
	}

	l.Infow("applicant.get_applicant.success")
	return &pb.GetApplicantResponse{Applicant: toPbApplicant(a)}, nil
}
func (s *service) GetApplicantByEmail(ctx context.Context, req *pb.GetApplicantByEmailRequest) (*pb.GetApplicantByEmailResponse, error) {
	l := s.log.With("op", "get_applicant_by_email", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "email", req.Email)

	if req.Email == "" {
		l.Errorw("applicant.get_applicant_by_email_failed", "err", "email cannot be empty")
		return nil, status.Errorf(codes.InvalidArgument, "email cannot be empty")
	}

	a, err := s.dataProvider.GetByEmail(ctx, req.Email)
	if err != nil {
		l.Errorw("applicant.get_applicant_by_email_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if a == nil {
		return nil, status.Errorf(codes.NotFound, "applicant not found")
	}

	l.Infow("applicant.get_applicant_by_email.success")
	return &pb.GetApplicantByEmailResponse{Applicant: toPbApplicant(a)}, nil
}

func (s *service) QueryApplicants(ctx context.Context, req *pb.QueryApplicantsRequest) (*pb.QueryApplicantsResponse, error) {
	l := s.log.With("op", "query_applicants", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "query", req)

	verr := validateQuery(req)
	if len(verr) > 0 {
		l.Errorw("applicant.query_applicants.validation_error", "err", verr)
		return nil, verr.ToStatus()
	}

	query := dal.NewQueryApplicantsDal(req.Ids, req.FullEmails, req.SubstrEmails,
		req.IsActive, req.IsDeleted, int(req.Page), int(req.PageSize),
	)
	list, err := s.dataProvider.QueryList(ctx, query)
	if err != nil {
		l.Errorw("applicant.query_applicants_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if len(list) == 0 {
		return nil, status.Errorf(codes.NotFound, "no applicants found")
	}

	result := make([]*pb.Applicant, 0, len(list))
	for _, a := range list {
		result = append(result, toPbApplicant(a))
	}

	l.Infow("applicant.query_applicants.success")
	return &pb.QueryApplicantsResponse{Applicants: result}, nil
}

func (s *service) createApplicant(req *pb.CreateApplicantRequest) (*applicant.Applicant, validationerror.ValidationError) {
	r := req.Applicant
	if r.Patronymic != nil && *r.Patronymic == "" {
		r.Patronymic = nil
	}
	if r.Contacts == nil {
		r.Contacts = &pb.Contacts{}
	}
	if r.Contacts.PhoneNumber != nil && *r.Contacts.PhoneNumber == "" {
		r.Contacts.PhoneNumber = nil
	}
	if r.Contacts.Telegram != nil && *r.Contacts.Telegram == "" {
		r.Contacts.Telegram = nil
	}

	return applicant.New(
		r.FirstName, r.LastName, r.Patronymic,
		r.BirthDate, r.City, r.Email,
		r.Contacts.PhoneNumber, r.Contacts.Telegram,
		false, false,
	)
}

func toPbApplicant(a *applicant.Applicant) *pb.Applicant {
	return &pb.Applicant{
		Id:         a.Id(),
		FirstName:  a.FirstName(),
		LastName:   a.LastName(),
		Patronymic: a.Patronymic(),
		BirthDate:  a.BirthDate(),
		City:       a.City(),
		Email:      a.Email(),
		Contacts: &pb.Contacts{
			PhoneNumber: a.PhoneNumber(),
			Telegram:    a.Telegram(),
		},
		IsActive:  a.IsActive(),
		IsDeleted: a.IsDeleted(),
		CreatedAt: timestamppb.New(a.CreatedAt()),
		UpdatedAt: timestamppb.New(a.UpdatedAt()),
	}
}

func validateQuery(req *pb.QueryApplicantsRequest) validationerror.ValidationError {
	verr := make(validationerror.ValidationError)

	if req.Page < 1 {
		verr["page"] = "page must be positive"
	}
	if req.PageSize < 1 {
		verr["page_size"] = "page_size must be positive"
	}

	for i, id := range req.Ids {
		if id < 1 {
			verr[fmt.Sprintf("ids[%d]", i)] = "id must be positive"
		}
	}
	for i, email := range req.FullEmails {
		if email == "" {
			verr[fmt.Sprintf("full_emails[%d]", i)] = "email cannot be empty"
		}
	}
	for i, emailSubstr := range req.SubstrEmails {
		if emailSubstr == "" {
			verr[fmt.Sprintf("substr_emails[%d]", i)] = "email cannot be empty"
		}
	}

	return verr
}
