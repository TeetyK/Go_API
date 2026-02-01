package main

import (
	"context"
	"log"
	"os"
	"time"

	"API/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultAddress = "localhost:50051"
)

func main() {
	// Check for command-line arguments
	if len(os.Args) < 4 {
		log.Fatalf("Usage: go run main.go <to> <subject> <body>")
	}
	to := os.Args[1]
	subject := os.Args[2]
	body := os.Args[3]

	// Set up a connection to the server.
	conn, err := grpc.Dial(defaultAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewNotificationServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Printf("Sending email to %s...", to)
	r, err := c.SendEmail(ctx, &proto.EmailRequest{To: to, Subject: subject, Body: body})
	if err != nil {
		log.Fatalf("could not send email: %v", err)
	}
	log.Printf("Server Response: %s", r.GetMessage())
}
