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

// create new user
func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	h.logger.Info().Msg("craete user")

	u := model.User{
		Username: req.User.Username,
		Email:    req.User.Email,
		Password: req.User.Password,
		Bio:      req.User.Bio,
		Image:    req.User.Image,
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

	return &pb.UserResponse{
		User: &pb.LoginedUser{
			Email:    u.Email,
			Token:    auth.GenerateJWT(u.ID),
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
		},
	}, nil
}

// login user
func (h *Handler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.UserResponse, error) {
	h.logger.Info().Msg("login user")

	u := model.User{}
	err := h.db.Where("email = ?", req.User.Email).First(&u).Error
	if err != nil {
		err = fmt.Errorf("failed to login due to wrong email: %w", err)
		h.logger.Error().Err(err)
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}

	if !u.CheckPassword(req.User.Password) {
		h.logger.Error().Msgf("failed to login due to receive wrong password: %s", u.Email)
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}

	return &pb.UserResponse{
		User: &pb.LoginedUser{
			Email:    u.Email,
			Token:    auth.GenerateJWT(u.ID),
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
		},
	}, nil
}
