package server

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func loggingMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	log.Println("Hello from loggingMiddleware!")

	resp, err := handler(ctx, req)

	log.Println("Bye from loggingMiddleware!")
	return resp, err
}
