package main

import (
	"context"
	"log"
	"net"
	"os"

	placesv1 "github.com/goghi48/ryden/gen/go/ryden/places/v1"
	"github.com/goghi48/ryden/internal/place/service"
	"github.com/goghi48/ryden/internal/place/storage"
	"github.com/goghi48/ryden/internal/place/transport"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	defer pool.Close()

	postgresStorage := storage.NewPostgresStorage(pool)

	placeService := service.NewPlaceService(postgresStorage)
	handler := transport.NewHandler(placeService)

	grpcServer := grpc.NewServer()

	placesv1.RegisterPlaceServiceServer(grpcServer, handler)

	log.Printf("place-service gRPC server started on :%s", grpcPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
