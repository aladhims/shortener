package shorten

import (
	"context"

	pb "github.com/aladhims/shortener/pkg/shorten/proto"
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

func (l *loggingService) Shorten(ctx context.Context, req *pb.ShortURL) (*pb.ShortenResponse, error) {
	log.WithFields(log.Fields{
		"origin": req.Origin,
		"slug":   req.Slug,
		"type":   req.UrlType,
	}).Info("A shorten request comes in")

	return l.service.Shorten(ctx, req)
}

func (l *loggingService) Expand(ctx context.Context, req *pb.ExpandRequest) (*pb.ExpandResponse, error) {
	log.WithFields(log.Fields{
		"slug": req.Slug,
	}).Info("An expand request comes in")

	return l.service.Expand(ctx, req)
}

func (l *loggingService) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {

	return l.service.Check(ctx, req)
}
