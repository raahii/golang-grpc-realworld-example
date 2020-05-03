package handler

import (
	"context"
	"fmt"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// show user profile
func (h *Handler) ShowProfile(ctx context.Context, req *pb.ShowProfileRequest) (*pb.ProfileResponse, error) {
	h.logger.Info().Msgf("Show profile | req: %+v\n", req)

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

	var u model.User
	err = h.db.Where("username = ?", req.Username).First(&u).Error
	if err != nil {
		h.logger.Fatal().Err(fmt.Errorf("user not found: %w", err))
		return nil, status.Error(codes.NotFound, "user was not found")
	}

	var count int
	err = h.db.Table("follows").Where("from_user_id = ? AND to_user_id = ?", currentUser.ID, u.ID).Count(&count).Error
	if err != nil {
		h.logger.Fatal().Err(err).Msg("failed to find following relationship")
		return nil, status.Error(codes.Aborted, "internal server error")
	}
	following := count >= 1

	p := pb.Profile{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: following,
	}

	return &pb.ProfileResponse{Profile: &p}, nil
}

// follow user
func (h *Handler) FollowUser(ctx context.Context, req *pb.FollowRequest) (*pb.ProfileResponse, error) {
	h.logger.Info().Msgf("Follow User | req: %+v\n", req)

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		err = fmt.Errorf("unauthenticated: %w", err)
		h.logger.Error().Err(err).Msg("unauthenticated")
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	var currentUser model.User
	err = h.db.Find(&currentUser, userID).Error
	if err != nil {
		h.logger.Fatal().Err(err).Msg("current user not found")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if currentUser.Username == req.Username {
		h.logger.Error().Msg("cannot follow yourself")
		return nil, status.Error(codes.InvalidArgument, "cannot follow yourself")
	}

	var u model.User
	err = h.db.Where("username = ?", req.Username).First(&u).Error
	if err != nil {
		h.logger.Fatal().Err(err).Msg("target user not found")
		return nil, status.Error(codes.NotFound, "user was not found")
	}

	err = h.db.Model(&currentUser).Association("Follows").Append(&u).Error
	if err != nil {
		h.logger.Fatal().Err(err).Msgf("failed to follow user: (ID: %d) -> (ID: %d)", currentUser.ID, u.ID)
		return nil, status.Error(codes.Aborted, "failed to follow user")
	}

	p := pb.Profile{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: true,
	}

	return &pb.ProfileResponse{Profile: &p}, nil
}
