package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/aladhims/shortener/pkg/notification"
	pb "github.com/aladhims/shortener/pkg/notification/proto"
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
	httpPort       string
	smtpHost       string
	smtpPort       string
	smtpUser       string
	smtpPass       string
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

	smtpConfig := notification.SmtpConfig{
		Host:     smtpHost,
		Port:     smtpPort,
		User:     smtpUser,
		Password: smtpPass,
	}

	srv := notification.NewService(smtpConfig)
	srv = notification.NewLoggingService(&logrus.Logger{
		Formatter: new(logrus.JSONFormatter),
		Level:     logrus.DebugLevel,
	}, srv)
	srv = notification.NewInstrumentingService(requestCounter, latencySummary, srv)

	pb.RegisterServiceServer(grpcServer, srv)

	grpcServer.Serve(lis)
}

func init() {
	os.Getenv("FOO")
	port = os.Getenv("PORT")
	host = os.Getenv("HOST")
	httpPort = os.Getenv("HTTP_PORT")
	serviceName = os.Getenv("SERVICE_NAME")
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")
	smtpUser = os.Getenv("SMTP_USER")
	smtpPass = os.Getenv("SMTP_PASSWORD")

	requestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "notify_service",
			Name:      "request_count",
			Help:      "Number of requests received",
		},
	)

	latencySummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "notify_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		},
	)
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(latencySummary)
}

func main() {
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

	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
