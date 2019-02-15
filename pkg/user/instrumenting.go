package user

import (
	"context"
	"time"

	pb "github.com/aladhims/shortener/pkg/user/proto"
	"github.com/prometheus/client_golang/prometheus"
)

type instrumentingService struct {
	requestCount   prometheus.Counter
	requestLatency prometheus.Histogram
	service        pb.ServiceServer
}

func NewInstrumentingService(counter prometheus.Counter, latency prometheus.Histogram, service pb.ServiceServer) pb.ServiceServer {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		service:        service,
	}
}

func (i *instrumentingService) Create(ctx context.Context, req *pb.User) (*pb.CreateResponse, error) {
	defer func(begin time.Time) {
		i.requestCount.Inc()
		i.requestLatency.Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.service.Create(ctx, req)
}

func (i *instrumentingService) Get(ctx context.Context, req *pb.GetRequest) (*pb.User, error) {
	defer func(begin time.Time) {
		i.requestCount.Inc()
		i.requestLatency.Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.service.Get(ctx, req)
}

func (i *instrumentingService) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	defer func(begin time.Time) {
		i.requestCount.Inc()
		i.requestLatency.Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.service.Check(ctx, req)
}
