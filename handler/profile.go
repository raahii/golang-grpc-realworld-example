package handler

import (
	"context"
	"fmt"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShowProfile gets a profile
func (h *Handler) ShowProfile(ctx context.Context, req *pb.ShowProfileRequest) (*pb.ProfileResponse, error) {
	h.logger.Info().Interface("req", req).Msg("show profile")

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("unauthenticated")
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	currentUser, err := h.us.GetByID(userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("current user not found")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	requestUser, err := h.us.GetByUsername(req.GetUsername())
	if err != nil {
		msg := "user was not found"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.NotFound, msg)
	}

	following, err := h.us.IsFollowing(currentUser, requestUser)
	if err != nil {
		msg := "failed to get following status"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.NotFound, "internal server error")
	}

	return &pb.ProfileResponse{Profile: requestUser.ProtoProfile(following)}, nil
}

// FollowUser follow a user
func (h *Handler) FollowUser(ctx context.Context, req *pb.FollowRequest) (*pb.ProfileResponse, error) {
	h.logger.Info().Interface("req", req).Msg("follow user")

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		err = fmt.Errorf("unauthenticated: %w", err)
		h.logger.Error().Err(err).Msg("unauthenticated")
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	currentUser, err := h.us.GetByID(userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("current user not found")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if currentUser.Username == req.GetUsername() {
		h.logger.Error().Msg("cannot follow yourself")
		return nil, status.Error(codes.InvalidArgument, "cannot follow yourself")
	}

	requestUser, err := h.us.GetByUsername(req.GetUsername())
	if err != nil {
		h.logger.Error().Err(fmt.Errorf("user not found: %w", err))
		return nil, status.Error(codes.NotFound, "user was not found")
	}

	err = h.us.Follow(currentUser, requestUser)
	if err != nil {
		msg := fmt.Sprintf("failed to follow user: (ID: %d) -> (ID: %d)",
			currentUser.ID, requestUser.ID)
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.Aborted, "failed to follow user")
	}

	return &pb.ProfileResponse{Profile: requestUser.ProtoProfile(true)}, nil
}

// UnfollowUser unfollow a user
func (h *Handler) UnfollowUser(ctx context.Context, req *pb.UnfollowRequest) (*pb.ProfileResponse, error) {
	h.logger.Info().Interface("req", req).Msg("unfollow user")

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		err = fmt.Errorf("unauthenticated: %w", err)
		h.logger.Error().Err(err).Msg("unauthenticated")
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	currentUser, err := h.us.GetByID(userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("current user not found")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if currentUser.Username == req.GetUsername() {
		h.logger.Error().Msg("cannot follow yourself")
		return nil, status.Error(codes.InvalidArgument, "cannot follow yourself")
	}

	requestUser, err := h.us.GetByUsername(req.GetUsername())
	if err != nil {
		msg := "user was not found"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.NotFound, msg)
	}

	following, err := h.us.IsFollowing(currentUser, requestUser)
	if err != nil {
		msg := "failed to get following status"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.NotFound, "internal server error")
	}

	if !following {
		h.logger.Error().Err(err).Msg("current user is not following request user")
		return nil, status.Errorf(codes.Unauthenticated, "you are not following the user")
	}

	err = h.us.Unfollow(currentUser, requestUser)
	if err != nil {
		msg := fmt.Sprintf("failed to unfollow user: (ID: %d) -> (ID: %d)",
			currentUser.ID, requestUser.ID)
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.Aborted, "failed to unfollow user")
	}

	return &pb.ProfileResponse{Profile: requestUser.ProtoProfile(false)}, nil
}
