package main

import (
	"context"
	"log"
	"net"
	"API/proto"
	"API/utils"
	"google.golang.org/grpc"
)

// server is used to implement proto.GreeterServer.
type server struct {
	proto.UnimplementedGreeterServer
}

// SayHello implements proto.GreeterServer
func (s *server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	log.Printf("Received: %v. Sending notification.", in.GetName())
	err := utils.SendPasswordResetEmail(in.GetName(), "grpc-token")
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return nil, err
	}
	return &proto.HelloReply{Message: "Hello " + in.GetName() + ", notification sent."}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
