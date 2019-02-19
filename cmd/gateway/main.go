package main

import (
	"log"
	"os"

	"github.com/aladhims/shortener/pkg/gateway"
	notificationpb "github.com/aladhims/shortener/pkg/notification/proto"
	shortenpb "github.com/aladhims/shortener/pkg/shorten/proto"
	userpb "github.com/aladhims/shortener/pkg/user/proto"
	"google.golang.org/grpc"
)

var (
	port string = "8082"

	shortenServiceHost string = "localhost"
	shortenServicePort string = "3032"
	shortenServiceAddr string
	shortenClient      shortenpb.ServiceClient

	userServiceHost string = "localhost"
	userServicePort string = "3033"
	userServiceAddr string
	userClient      userpb.ServiceClient

	notificationServiceHost string
	notificationServicePort string = "3034"
	notificationServiceAddr string
	notificationClient      notificationpb.ServiceClient
)

func init() {
	port = os.Getenv("PORT")

	shortenServiceHost = os.Getenv("SHORTEN_HOST")
	shortenServicePort = os.Getenv("SHORTEN_PORT")

	shortenServiceAddr = shortenServiceHost + ":" + shortenServicePort

	userServiceHost = os.Getenv("USER_HOST")
	userServicePort = os.Getenv("USER_PORT")

	userServiceAddr = userServiceHost + ":" + userServicePort

	notificationServiceHost = os.Getenv("NOTIFICATION_HOST")
	notificationServicePort = os.Getenv("NOTIFICATION_PORT")

	notificationServiceAddr = notificationServiceHost + ":" + notificationServicePort
}

func main() {
	shortenConn, err := grpc.Dial(shortenServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't connect")
	}
	defer shortenConn.Close()

	shortenClient = shortenpb.NewServiceClient(shortenConn)

	userConn, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't connect")
	}

	defer userConn.Close()

	userClient = userpb.NewServiceClient(userConn)

	notificationConn, err := grpc.Dial(notificationServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't connect")
	}

	defer notificationConn.Close()

	notificationClient = notificationpb.NewServiceClient(notificationConn)

	frontendService := gateway.NewService(shortenClient, userClient, notificationClient)
	frontendService.Run(port)
}
