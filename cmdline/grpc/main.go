package main

import (
	"context"
	"log"
	"net"
	"time"

	"API/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	proto.UnimplementedNotificationServiceServer
}

func (s *server) SendEmail(ctx context.Context, in *proto.EmailRequest) (*proto.EmailReply, error) {
	log.Printf("gRPC proxy: Received SendEmail request for: %v. Forwarding...", in.GetTo())
	downstreamConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("gRPC proxy: Failed to connect to downstream service: %v", err)
		return nil, err
	}
	defer downstreamConn.Close()
	downstreamClient := proto.NewNotificationServiceClient(downstreamConn)
	downstreamCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log.Printf("gRPC proxy: Forwarding request to downstream service at %s", downstreamConn.Target())
	reply, err := downstreamClient.SendEmail(downstreamCtx, in)
	if err != nil {
		log.Printf("gRPC proxy: Error from downstream service: %v", err)
		return nil, err
	}
	log.Printf("gRPC proxy: Successfully received reply from downstream: '%s'", reply.GetMessage())

	return reply, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterNotificationServiceServer(s, &server{})
	log.Printf("gRPC server (proxy) listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
