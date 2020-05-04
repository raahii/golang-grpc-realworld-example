package handler

import (
	"context"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateArticle creates a article
func (h *Handler) CreateArticle(ctx context.Context, req *pb.CreateAritcleRequest) (*pb.ArticleResponse, error) {
	h.logger.Info().Msgf("Create artcile | req: %+v\n", req)

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("unauthenticated")
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	var currentUser model.User
	err = h.db.Find(&currentUser, userID).Error
	if err != nil {
		h.logger.Fatal().Err(err).Msg("current user not found")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.ArticleResponse{}, nil
}
