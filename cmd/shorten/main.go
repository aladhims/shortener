package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/aladhims/shortener/pkg/shorten"
	pb "github.com/aladhims/shortener/pkg/shorten/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	selfClient     pb.ServiceClient
	mongoClient    *mongo.Client
	mongoHost      string = "localhost"
	mongoPort      string = "27017"
	selfAddr       string
	serviceName    string        = "shorten"
	host           string        = "localhost"
	port           string        = "3032"
	httpPort       string        = "3042"
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
	mongoClient, err := mongo.NewClient(fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort))
	if err != nil {
		log.Fatalf("failed to connect to mongoDB: %v", err)
	}

	repo := shorten.NewMongoRepository(mongoClient)
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
	httpPort = os.Getenv("HTTP_PORT")
	port = os.Getenv("PORT")
	host = os.Getenv("HOST")
	mongoHost = os.Getenv("MONGO_HOST")
	mongoPort = os.Getenv("MONGO_PORT")
	serviceName = os.Getenv("SERVICE_NAME")

	selfAddr = host + ":" + port

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

	go runGRPCServer()

	selfConn, err := grpc.Dial(selfAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't connect")
	}

	defer selfConn.Close()

	selfClient = pb.NewServiceClient(selfConn)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", handleHealthCheck)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		_, err := selfClient.Shorten(context.Background(), &pb.ShortURL{
			UserId:  121212,
			Origin:  "https://www.aladhims.com",
			Slug:    "me",
			UrlType: pb.URLType_DEFINED,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "OK")
	})

	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
