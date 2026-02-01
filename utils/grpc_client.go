package utils

import (
	"API/proto"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcAddress = "localhost:50051"
)

// SendEmailViaGRPC connects to the gRPC server and calls the SendEmail RPC.
func SendEmailViaGRPC(to, subject, body string) (string, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect to gRPC server: %v", err)
		return "", fmt.Errorf("did not connect to gRPC server: %w", err)
	}
	defer conn.Close()
	c := proto.NewNotificationServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Sending gRPC email request to %s...", to)
	r, err := c.SendEmail(ctx, &proto.EmailRequest{To: to, Subject: subject, Body: body})
	if err != nil {
		log.Printf("gRPC call failed: %v", err)
		return "", fmt.Errorf("gRPC call failed: %w", err)
	}

	log.Printf("gRPC Server Response: %s", r.GetMessage())
	return r.GetMessage(), nil
}
