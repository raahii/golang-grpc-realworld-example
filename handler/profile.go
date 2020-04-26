package handler

import (
	"context"
	"fmt"

	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// show user profile
func (h *Handler) ShowProfile(ctx context.Context, req *pb.ShowProfileRequest) (*pb.ProfileResponse, error) {
	h.logger.Info().Msgf("Show profile | req: %+v\n", req)

	user := model.User{}
	err := h.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		h.logger.Fatal().Err(fmt.Errorf("user not found: %w", err))
		return nil, status.Error(codes.NotFound, "user was not found")
	}

	p := pb.Profile{
		Username: user.Username,
		Bio:      user.Bio,
		Image:    user.Image,
	}

	return &pb.ProfileResponse{Profile: &p}, nil
}
