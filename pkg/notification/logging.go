package notification

import (
	"context"

	pb "github.com/aladhims/shortener/pkg/notification/proto"
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

func (l *loggingService) Notify(ctx context.Context, req *pb.NotifyRequest) (*pb.NotifyResponse, error) {

	log.WithFields(log.Fields{
		"email":    req.Email,
		"fullname": req.Fullname,
		"slug":     req.Slug,
		"origin":   req.Origin,
	}).Info("A notify request comes in")

	return l.service.Notify(ctx, req)
}

func (l *loggingService) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {

	return l.service.Check(ctx, req)
}
