package user

import (
	"context"

	pb "github.com/aladhims/shortener/pkg/user/proto"
	log "github.com/sirupsen/logrus"
)

type loggingService struct {
	logger  *log.Logger
	service pb.ServiceServer
}

func NewLoggingService(logger *log.Logger, service pb.ServiceServer) pb.ServiceServer {
	return &loggingService{
		logger:  logger,
		service: service,
	}
}

func (l *loggingService) Create(ctx context.Context, req *pb.User) (*pb.CreateResponse, error) {
	log.WithFields(log.Fields{
		"fullname": req.Fullname,
		"email":    req.Email,
	}).Info("A create user request comes in")

	return l.service.Create(ctx, req)
}

func (l *loggingService) Get(ctx context.Context, req *pb.GetRequest) (*pb.User, error) {
	log.WithFields(log.Fields{
		"ID": req.Id,
	}).Info("A get user request comes in")

	return l.service.Get(ctx, req)
}

func (l *loggingService) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {

	return l.service.Check(ctx, req)
}
