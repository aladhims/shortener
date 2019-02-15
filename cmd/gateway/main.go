package main

import (
	"log"

	"github.com/aladhims/shortener/pkg/gateway"
	notificationpb "github.com/aladhims/shortener/pkg/notification/proto"
	shortenpb "github.com/aladhims/shortener/pkg/shorten/proto"
	userpb "github.com/aladhims/shortener/pkg/user/proto"
	"google.golang.org/grpc"
)

var (
	shortenServiceHost string
	shortenServicePort string
	shortenServiceAddr string
	shortenClient      shortenpb.ServiceClient

	userServiceHost string
	userServicePort string
	userServiceAddr string
	userClient      userpb.ServiceClient

	notificationServiceHost string
	notificationServicePort string
	notificationServiceAddr string
	notificationClient      notificationpb.ServiceClient
)

func main() {

	shortenServiceHost = "shorten-service"
	shortenServicePort = "3032"

	shortenServiceAddr = shortenServiceHost + ":" + shortenServicePort
	shortenConn, err := grpc.Dial(shortenServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't connect")
	}

	defer shortenConn.Close()

	shortenClient = shortenpb.NewServiceClient(shortenConn)

	userServiceHost = "user-service"
	userServicePort = "3033"

	userServiceAddr = userServiceHost + ":" + userServicePort

	userConn, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't connect")
	}

	defer userConn.Close()

	userClient = userpb.NewServiceClient(userConn)

	notificationServiceHost = "notification-service"
	notificationServicePort = "3034"

	notificationServiceAddr = notificationServiceHost + ":" + notificationServicePort

	notificationConn, err := grpc.Dial(notificationServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can't connect")
	}

	defer notificationConn.Close()

	notificationClient = notificationpb.NewServiceClient(notificationConn)

	frontendService := gateway.NewService(shortenClient, userClient, notificationClient)
	frontendService.Run()
}
