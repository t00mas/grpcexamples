package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	service_v1 "github.com/t00mas/grpcexamples/proto/gen/go/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	defaultMsg = "keepalive demo"
)

var (
	app  = flag.String("app", "server", "which app to execute (server / client)")
	host = flag.String("host", "localhost", "the host to connect to")
	msg  = flag.String("msg", defaultMsg, "message to send")
	port = flag.Int("port", 50051, "The server port")
	kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}
	kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second,
		MaxConnectionAge:      30 * time.Second,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  5 * time.Second,
		Timeout:               1 * time.Second,
	}
	kacp = keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}
)

type EchoService struct {
	service_v1.UnimplementedServiceServer
}

func (s *EchoService) Echo(ctx context.Context, in *service_v1.Request) (*service_v1.Response, error) {
	log.Printf("Message received: %v", in.Message)
	return &service_v1.Response{Message: in.Message}, nil
}

func main() {
	flag.Parse()
	switch *app {
	case "server":
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
		service_v1.RegisterServiceServer(s, &EchoService{})

		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	case "client":
		addr := fmt.Sprintf("%s:%d", *host, *port)

		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		c := service_v1.NewServiceClient(conn)
		r, err := c.Echo(ctx, &service_v1.Request{Message: *msg})
		if err != nil {
			log.Fatalf("could not send message: %v", err)
		}

		log.Printf("RPC response: %s", r.GetMessage())
		select {}
	}
}
