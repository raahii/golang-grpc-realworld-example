package handler

import (
	"context"
	"fmt"
	"strconv"

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

// GetArticle gets a article
func (h *Handler) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.ArticleResponse, error) {
	// get article
	articleID, err := strconv.Atoi(req.GetSlug())
	if err != nil {
		msg := fmt.Sprintf("cannot convert slug (%s) into integer", req.GetSlug())
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.InvalidArgument, "invalid article id")
	}

	var a model.Article
	err = h.db.Preload("Tags").Find(&a, articleID).Error
	if err != nil {
		msg := fmt.Sprintf("requested article (slug=%d) not found", articleID)
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.InvalidArgument, "invalid article id")
	}

	// get author
	var u model.User
	err = h.db.Find(&u, a.UserID).Error
	if err != nil {
		msg := "article author not found"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.NotFound, msg)
	}

	// bind article
	var pa pb.Article
	err = a.BindTo(&pa)
	if err != nil {
		msg := "failed to convert model.User to pb.User"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.Aborted, "internal server error")
	}

	pa.Author = &pb.Profile{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: false,
	}

	// get current user if exists
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		pa.Favorited = false
		pa.Author.Following = false
		return &pb.ArticleResponse{Article: &pa}, nil
	}

	var currentUser model.User
	err = h.db.Find(&currentUser, userID).Error
	if err != nil {
		msg := fmt.Sprintf("token is valid but the user not found")
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.NotFound, msg)
	}

	// TODO: set favorite field
	pa.Favorited = false

	// set following field
	var count int
	err = h.db.Table("follows").Where("from_user_id = ? AND to_user_id = ?", currentUser.ID, a.UserID).Count(&count).Error
	if err != nil {
		h.logger.Fatal().Err(err).Msg("failed to find following relationship")
		return nil, status.Error(codes.Aborted, "internal server error")
	}
	pa.Author.Following = count >= 1

	return &pb.ArticleResponse{Article: &pa}, nil
}
