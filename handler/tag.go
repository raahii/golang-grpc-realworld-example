package handler

import (
	"context"

	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetTags returns all of tags
func (h *Handler) GetTags(ctx context.Context, req *pb.Empty) (*pb.TagsResponse, error) {
	h.logger.Info().Interface("req", req).Msg("get tags")

	tags, err := h.as.GetTags()
	if err != nil {
		h.logger.Error().Err(err).Msg("faield to get tags")
		return nil, status.Error(codes.Aborted, "internal server error")
	}

	tagNames := make([]string, 0, len(tags))
	for _, t := range tags {
		tagNames = append(tagNames, t.Name)
	}

	return &pb.TagsResponse{Tags: tagNames}, nil
}
