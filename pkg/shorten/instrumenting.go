package shorten

import (
	"context"
	"time"

	pb "github.com/aladhims/shortener/pkg/shorten/proto"
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

func (i *instrumentingService) Shorten(ctx context.Context, req *pb.ShortURL) (*pb.ShortenResponse, error) {
	defer func(begin time.Time) {
		i.requestCount.Inc()
		i.requestLatency.Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.service.Shorten(ctx, req)
}

func (i *instrumentingService) Expand(ctx context.Context, req *pb.ExpandRequest) (*pb.ExpandResponse, error) {
	defer func(begin time.Time) {
		i.requestCount.Inc()
		i.requestLatency.Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.service.Expand(ctx, req)
}

func (i *instrumentingService) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	defer func(begin time.Time) {
		i.requestCount.Inc()
		i.requestLatency.Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.service.Check(ctx, req)
}
