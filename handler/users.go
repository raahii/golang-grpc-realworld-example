package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/k0kubun/pp"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

func (h *Handler) ShowProfile(ctx context.Context, req *pb.ShowProfileRequest) (*pb.ShowProfileResponse, error) {
	h.logger.Printf("Show profile | req: %+v\n", req)

	user := model.User{}
	err := h.db.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		log.Fatal(fmt.Errorf("user not found: %w", err))
	}

	pp.Println(user)

	var bio string
	if user.Bio != nil {
		bio = *user.Bio
	}

	var image string
	if user.Image != nil {
		image = *user.Image
	}

	p := pb.Profile{
		Username: req.Username,
		Bio:      bio,
		Image:    image,
	}

	return &pb.ShowProfileResponse{Profile: &p}, nil
}
