package employerservice

import (
	"context"
	"fmt"

	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	pb "github.com/ZaiiiRan/job_search_service/user-service/gen/go/user_service/v1"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/domain/user/employer"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/errors/validationerror"
	dal "github.com/ZaiiiRan/job_search_service/user-service/internal/repositories/models"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/transport/postgres"
	"github.com/ZaiiiRan/job_search_service/user-service/internal/transport/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EmployerService interface {
	CreateEmployer(ctx context.Context, req *pb.CreateEmployerRequest) (*pb.CreateEmployerResponse, error)
	GetEmployer(ctx context.Context, req *pb.GetEmployerRequest) (*pb.GetEmployerResponse, error)
	GetEmployerByEmail(ctx context.Context, req *pb.GetEmployerByEmailRequest) (*pb.GetEmployerByEmailResponse, error)
	QueryEmployers(ctx context.Context, req *pb.QueryEmployersRequest) (*pb.QueryEmployersResponse, error)
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *employerDataProvider
}

func New(pgClient *postgres.PostgresClient, redisClient *redis.RedisClient, log *zap.SugaredLogger) EmployerService {
	return &service{log: log, dataProvider: newEmployerDataProvider(pgClient, redisClient)}
}

func (s *service) CreateEmployer(ctx context.Context, req *pb.CreateEmployerRequest) (*pb.CreateEmployerResponse, error) {
	l := s.log.With("op", "create_employer", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	e, verr := s.createEmployer(req.Employer)
	if len(verr) > 0 {
		l.Errorw("employer.create_employer_failed.validation_error", "err", verr)
		return nil, verr.ToStatus()
	}

	existed, err := s.dataProvider.GetByEmail(ctx, e.Email())
	if err != nil {
		l.Errorw("employer.create_employer_failed.check_existing_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if existed != nil {
		if existed.IsActive() {
			l.Errorw("employer.create_employer_failed", "err", "employer with this email already exists")
			return nil, status.Errorf(codes.AlreadyExists, "employer with this email already exists")
		}
		l.Infow("employer.create_employer.restoring_inactive_employer", "id", existed.Id())
		e.SetId(existed.Id())
	}

	if err := s.dataProvider.Save(ctx, e); err != nil {
		l.Errorw("employer.create_employer_failed.save_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	l.Infow("employer.create_employer.created")
	return &pb.CreateEmployerResponse{Employer: toPbEmployer(e)}, nil
}

func (s *service) GetEmployer(ctx context.Context, req *pb.GetEmployerRequest) (*pb.GetEmployerResponse, error) {
	l := s.log.With("op", "get_employer", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "id", req.Id)

	if req.Id < 1 {
		l.Errorw("employer.get_employer_failed", "err", "id must be positive")
		return nil, status.Errorf(codes.InvalidArgument, "id must be positive")
	}

	e, err := s.dataProvider.GetById(ctx, req.Id)
	if err != nil {
		l.Errorw("employer.get_employer_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if e == nil {
		return nil, status.Errorf(codes.NotFound, "employer not found")
	}

	l.Infow("employer.get_employer.success")
	return &pb.GetEmployerResponse{Employer: toPbEmployer(e)}, nil
}

func (s *service) GetEmployerByEmail(ctx context.Context, req *pb.GetEmployerByEmailRequest) (*pb.GetEmployerByEmailResponse, error) {
	l := s.log.With("op", "get_employer_by_email", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "email", req.Email)

	if req.Email == "" {
		l.Errorw("employer.get_employer_by_email_failed", "err", "email cannot be empty")
		return nil, status.Errorf(codes.InvalidArgument, "email cannot be empty")
	}

	e, err := s.dataProvider.GetByEmail(ctx, req.Email)
	if err != nil {
		l.Errorw("employer.get_employer_by_email_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if e == nil {
		return nil, status.Errorf(codes.NotFound, "employer not found")
	}

	l.Infow("employer.get_employer_by_email.success")
	return &pb.GetEmployerByEmailResponse{Employer: toPbEmployer(e)}, nil
}

func (s *service) QueryEmployers(ctx context.Context, req *pb.QueryEmployersRequest) (*pb.QueryEmployersResponse, error) {
	l := s.log.With("op", "query_employers", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "query", req)

	verr := validateQuery(req)
	if len(verr) > 0 {
		l.Errorw("employer.query_employers.validation_error", "err", verr)
		return nil, verr.ToStatus()
	}

	query := dal.NewQueryEmployersDal(req.Ids, req.FullEmails, req.SubstrEmails, req.FullCompanyNames,
		req.SubstrCompanyNames, req.IsActive, req.IsDeleted, int(req.Page), int(req.PageSize),
	)
	list, err := s.dataProvider.QueryList(ctx, query)
	if err != nil {
		l.Errorw("employer.query_employers_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if len(list) == 0 {
		return nil, status.Errorf(codes.NotFound, "no employers found")
	}

	result := make([]*pb.Employer, 0, len(list))
	for _, e := range list {
		result = append(result, toPbEmployer(e))
	}

	l.Infow("employer.query_employers.success")
	return &pb.QueryEmployersResponse{Employers: result}, nil
}

func (s *service) createEmployer(r *pb.Employer) (*employer.Employer, validationerror.ValidationError) {
	if r.Contacts == nil {
		r.Contacts = &pb.Contacts{}
	}
	if r.Contacts.PhoneNumber != nil && *r.Contacts.PhoneNumber == "" {
		r.Contacts.PhoneNumber = nil
	}
	if r.Contacts.Telegram != nil && *r.Contacts.Telegram == "" {
		r.Contacts.Telegram = nil
	}

	return employer.New(
		r.CompanyName, r.City, r.Email,
		r.Contacts.PhoneNumber, r.Contacts.Telegram,
		false, false,
	)
}

func toPbEmployer(e *employer.Employer) *pb.Employer {
	return &pb.Employer{
		Id:          e.Id(),
		CompanyName: e.CompanyName(),
		City:        e.City(),
		Email:       e.Email(),
		Contacts: &pb.Contacts{
			PhoneNumber: e.PhoneNumber(),
			Telegram:    e.Telegram(),
		},
		IsActive:  e.IsActive(),
		IsDeleted: e.IsDeleted(),
		CreatedAt: timestamppb.New(e.CreatedAt()),
		UpdatedAt: timestamppb.New(e.UpdatedAt()),
	}
}

func validateQuery(req *pb.QueryEmployersRequest) validationerror.ValidationError {
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
	for i, companyName := range req.FullCompanyNames {
		if companyName == "" {
			verr[fmt.Sprintf("full_company_names[%d]", i)] = "email cannot be empty"
		}
	}
	for i, companyNameSubstr := range req.SubstrCompanyNames {
		if companyNameSubstr == "" {
			verr[fmt.Sprintf("substrs_company_names[%d]", i)] = "email cannot be empty"
		}
	}

	return verr
}
