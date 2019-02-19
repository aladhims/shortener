package user

import (
	"context"

	pb "github.com/aladhims/shortener/pkg/user/proto"
)

type service struct {
	repository Repository
}

// NewService creates a new service with the given repository
func NewService(repository Repository) pb.ServiceServer {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(ctx context.Context, req *pb.User) (*pb.CreateResponse, error) {

	existingUser := s.repository.GetByEmail(req.Email)
	if existingUser != nil {
		return &pb.CreateResponse{
			Id: existingUser.Id,
		}, nil
	}

	id, err := s.repository.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{
		Id: id,
	}, nil
}

func (s *service) Get(ctx context.Context, req *pb.GetRequest) (*pb.User, error) {
	user := s.repository.Get(req.Id)

	return user, nil
}

func (s *service) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	if s.repository == nil {
		return &pb.HealthCheckResponse{
			Status: pb.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	if err := s.repository.Ping(); err != nil {
		return &pb.HealthCheckResponse{
			Status: pb.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}
