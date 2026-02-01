package main

import (
	"context"
	"log"
	"net"
	"os"

	"API/proto"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedNotificationServiceServer
}

func (s *server) SendEmail(ctx context.Context, in *proto.EmailRequest) (*proto.EmailReply, error) {
	log.Printf("FINAL SERVER: Received SendEmail request for: %v. Subject: %s", in.GetTo(), in.GetSubject())

	// This is the final server, so we just simulate the email sending.
	log.Printf("FINAL SERVER: Simulating successful email send to %s", in.GetTo())

	return &proto.EmailReply{Message: "Email sent successfully from final server"}, nil
}

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052" // Default to 50051 for the final server
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterNotificationServiceServer(s, &server{})
	log.Printf("gRPC final server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
