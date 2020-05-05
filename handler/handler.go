package handler

import (
	"context"

	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"github.com/rs/zerolog"
)

// Handler definition
type Handler struct {
	logger *zerolog.Logger
	db     *gorm.DB
	us     *store.UserStore
	as     *store.ArticleStore
}

// New returns a new handler with logger and database
func New(l *zerolog.Logger, d *gorm.DB, us *store.UserStore, as *store.ArticleStore) *Handler {
	return &Handler{logger: l, db: d, us: us, as: as}
}

// SayHello is a dummy method
func (h *Handler) SayHello(ctx context.Context, in *pb.Empty) (*pb.HelloReply, error) {
	h.logger.Info().Msgf("hello request")
	return &pb.HelloReply{Message: "Hello, World!"}, nil
}
