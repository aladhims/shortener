package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/aladhims/shortener/pkg/shorten"
	pb "github.com/aladhims/shortener/pkg/shorten/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	selfClient     pb.ServiceClient
	selfAddr       string
	serviceName    string
	host           string
	port           string
	timeoutDur     time.Duration = time.Second
	latencySummary prometheus.Summary
	requestCounter prometheus.Counter
)

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if selfClient == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "NOT_HEALTHY")
		logrus.Errorf("server is not healthy, unable to connect remote grpc server")
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), timeoutDur)
	resp, err := selfClient.Check(ctx, &pb.HealthCheckRequest{Service: serviceName})
	if err == nil && resp.Status == pb.HealthCheckResponse_SERVING {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
		logrus.Infof("health check is OK")
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "NOT_HEALTHY")
	logrus.Errorf("server is not healthy err=%v response=%v", err, resp)
}

func runGRPCServer() {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	repo := shorten.NewInmemRepository()
	srv := shorten.NewService(repo)
	srv = shorten.NewLoggingService(&logrus.Logger{
		Formatter: new(logrus.JSONFormatter),
		Level:     logrus.DebugLevel,
	}, srv)
	srv = shorten.NewInstrumentingService(requestCounter, latencySummary, srv)

	pb.RegisterServiceServer(grpcServer, srv)

	grpcServer.Serve(lis)
}

func init() {
	requestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "shorten_service",
			Name:      "request_count",
			Help:      "Number of requests received",
		},
	)

	latencySummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "shorten_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		},
	)
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(latencySummary)
}

func main() {
	port = "3032"
	host = "shorten-service"
	serviceName = "shorten"

	go runGRPCServer()

	selfAddr = host + ":" + port
	selfConn, err := grpc.Dial(selfAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't connect")
	}

	defer selfConn.Close()

	selfClient = pb.NewServiceClient(selfConn)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", handleHealthCheck)

	log.Fatal(http.ListenAndServe(":3042", nil))
}
