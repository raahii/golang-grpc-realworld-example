package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/db"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type Logger interface {
	Printf(string, ...interface{})
}

type server struct {
	logger Logger
	db     *gorm.DB
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	d, err := db.New()
	if err != nil {
		l.Fatal(fmt.Errorf("failed to connect database: %w", err))
	}
	db.AutoMigrate(d)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		l.Fatal(fmt.Errorf("failed to listen: %w", err))
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{logger: l, db: d})
	l.Printf("starting server on port %s\n", port)
	if err := s.Serve(lis); err != nil {
		l.Fatal(fmt.Errorf("failed to serve: %w", err))
	}
}
