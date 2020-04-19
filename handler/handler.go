package handler

import (
	"context"
	"log"

	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

type Logger interface {
	Printf(string, ...interface{})
	Fatal(...interface{})
}

type Handler struct {
	logger Logger
	db     *gorm.DB
}

func New(l Logger, d *gorm.DB) *Handler {
	return &Handler{logger: l, db: d}
}

func (s *Handler) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
