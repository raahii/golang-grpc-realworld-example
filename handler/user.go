package handler

import (
	"context"
	"fmt"

	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ShowProfile(ctx context.Context, req *pb.ShowProfileRequest) (*pb.ShowProfileResponse, error) {
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

	return &pb.ShowProfileResponse{Profile: &p}, nil
}

func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
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
		msg := fmt.Sprintf("Validation error. %s", err)
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	err = u.HashPassword()
	if err != nil {
		return nil, status.Error(codes.Aborted, "Internal sever error.")
	}

	err = h.db.Create(&u).Error
	if err != nil {
		msg := fmt.Sprintf("Failed to create user. %s", err)
		return nil, status.Error(codes.Canceled, msg)
	}

	return &pb.CreateUserResponse{
		User: &pb.CreatedUser{
			Email:    u.Email,
			Token:    "", // TODO
			Username: u.Username,
			Bio:      u.Bio,
			Image:    u.Image,
		},
	}, nil
}
