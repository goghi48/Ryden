package main

import (
	"log"
	"net"

	placesv1 "github.com/goghi48/ryden/gen/go/ryden/places/v1"
	"github.com/goghi48/ryden/internal/place/service"
	"github.com/goghi48/ryden/internal/place/storage"
	"github.com/goghi48/ryden/internal/place/transport"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	memoryStorage := storage.NewMemoryStorage()
	placeService := service.NewPlaceService(memoryStorage)
	handler := transport.NewHandler(placeService)

	grpcServer := grpc.NewServer()

	placesv1.RegisterPlaceServiceServer(grpcServer, handler)

	log.Println("place-service gRPC server started on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
