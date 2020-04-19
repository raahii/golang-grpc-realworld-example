package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/raahii/golang-grpc-realworld-example/db"
	"github.com/raahii/golang-grpc-realworld-example/handler"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	d, err := db.New()
	if err != nil {
		l.Fatal(fmt.Errorf("failed to connect database: %w", err))
	}
	db.AutoMigrate(d)

	h := handler.New(l, d)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		l.Fatal(fmt.Errorf("failed to listen: %w", err))
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, h)
	pb.RegisterUsersServer(s, h)
	l.Printf("starting server on port %s\n", port)
	if err := s.Serve(lis); err != nil {
		l.Fatal(fmt.Errorf("failed to serve: %w", err))
	}
}
