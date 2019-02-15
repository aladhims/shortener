package notification

import (
	"context"
	"fmt"

	pb "github.com/aladhims/shortener/pkg/notification/proto"
)

const sender = "shortener@aladhims.com"

type service struct {
	smtpConfig SmtpConfig
}

func NewService(c SmtpConfig) pb.ServiceServer {
	return &service{
		smtpConfig: c,
	}
}

func (s *service) Notify(ctx context.Context, req *pb.NotifyRequest) (*pb.NotifyResponse, error) {

	body := fmt.Sprintf("Hello %s, the slug for %s is %s", req.Fullname, req.Origin, req.Slug)

	mail := mail{
		sender:  sender,
		to:      []string{req.Email},
		subject: "shorten request",
		cc:      []string{"shorten request detail"},
		body:    body,
	}

	err := s.smtpConfig.sendMail(mail)
	if err != nil {
		return nil, err
	}

	return &pb.NotifyResponse{}, nil

}

func (s *service) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	// if err := s.smtpClient.Noop(); err != nil {
	// 	return &pb.HealthCheckResponse{
	// 		Status: pb.HealthCheckResponse_NOT_SERVING,
	// 	}, nil
	// }

	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}
