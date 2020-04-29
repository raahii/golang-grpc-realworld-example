package handler

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// login user
func (h *Handler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.UserResponse, error) {
	h.logger.Info().Msg("login user")

	u := model.User{}
	err := h.db.Where("email = ?", req.User.GetEmail()).First(&u).Error
	if err != nil {
		err = fmt.Errorf("failed to login due to wrong email: %w", err)
		h.logger.Error().Err(err)
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}

	if !u.CheckPassword(req.User.GetPassword()) {
		h.logger.Error().Msgf("failed to login due to receive wrong password: %s", u.Email)
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		err := fmt.Errorf("Failed to create token. %w", err)
		h.logger.Error().Err(err)
		return nil, status.Error(codes.Aborted, "internal server error")
	}

	return &pb.UserResponse{
		User: &pb.User{
			Email:    u.Email,
			Token:    token,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
		},
	}, nil
}

// create new user
func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	h.logger.Info().Msg("craete user")

	u := model.User{
		Username: req.User.GetUsername(),
		Email:    req.User.GetEmail(),
		Password: req.User.GetPassword(),
	}

	err := u.Validate()
	if err != nil {
		err = fmt.Errorf("validation error: %w", err)
		h.logger.Error().Err(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = u.HashPassword()
	if err != nil {
		err := fmt.Errorf("Failed to hash password, %w", err)
		h.logger.Error().Err(err)
		return nil, status.Error(codes.Aborted, err.Error())
	}

	err = h.db.Create(&u).Error
	if err != nil {
		err := fmt.Errorf("Failed to create user. %w", err)
		h.logger.Error().Err(err)
		return nil, status.Error(codes.Canceled, err.Error())
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		err := fmt.Errorf("Failed to create token. %w", err)
		h.logger.Error().Err(err)
		return nil, status.Error(codes.Aborted, "internal server error")
	}

	return &pb.UserResponse{
		User: &pb.User{
			Email:    u.Email,
			Token:    token,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
		},
	}, nil
}

// get current user
func (h *Handler) CurrentUser(ctx context.Context, req *empty.Empty) (*pb.UserResponse, error) {
	h.logger.Info().Msg("get current user")

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		h.logger.Error().Err(err)
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	u := model.User{}
	err = h.db.Where("id = ?", userID).First(&u).Error
	if err != nil {
		err = fmt.Errorf("token is valid but the user not found: %w", err)
		h.logger.Error().Err(err)
		return nil, status.Error(codes.NotFound, "not user found")
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		err := fmt.Errorf("Failed to create token. %w", err)
		h.logger.Error().Err(err)
		return nil, status.Error(codes.Aborted, "internal server error")
	}

	return &pb.UserResponse{
		User: &pb.User{
			Email:    u.Email,
			Token:    token,
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
		},
	}, nil
}
