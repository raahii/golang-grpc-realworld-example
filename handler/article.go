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

	ra := req.GetArticle()
	tags := make([]model.Tag, 0, len(ra.GetTagList()))
	for _, t := range ra.GetTagList() {
		tags = append(tags, model.Tag{Name: t})
	}

	a := model.Article{
		Title:       ra.GetTitle(),
		Description: ra.GetDescription(),
		Body:        ra.GetBody(),
		UserID:      currentUser.ID,
		Tags:        tags,
	}

	err = a.Validate()
	if err != nil {
		msg := "validation error"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	err = h.db.Create(&a).Error
	if err != nil {
		msg := "Failed to create user."
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.Canceled, msg)
	}

	var pa pb.Article
	err = a.BindTo(&pa)
	if err != nil {
		msg := "Failed to convert model.User to pb.User"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.Aborted, "internal server error")
	}

	pa.Favorited = false

	pa.Author = &pb.Profile{
		Username:  currentUser.Username,
		Bio:       currentUser.Bio,
		Image:     currentUser.Image,
		Following: false,
	}

	return &pb.ArticleResponse{Article: &pa}, nil
}
