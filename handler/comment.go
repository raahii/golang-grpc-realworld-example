package handler

import (
	"context"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

// CreateComment create a comment for an article
func (h *Handler) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CommentResponse, error) {
	return &pb.CommentResponse{Comment: &pb.Comment{}}, nil
}
