package handler

import (
	"context"

	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/rs/zerolog"
)

type Handler struct {
	logger *zerolog.Logger
	db     *gorm.DB
}

func New(l *zerolog.Logger, d *gorm.DB) *Handler {
	return &Handler{logger: l, db: d}
}

func (h *Handler) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	h.logger.Info().Msgf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
