package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/k0kubun/pp"
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
		Author:      currentUser,
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
	err = a.BindTo(&pa, &currentUser, h.db)
	if err != nil {
		msg := "Failed to convert model.User to pb.User"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.Aborted, "internal server error")
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
	err = h.db.Preload("Tags").Preload("Author").Find(&a, articleID).Error
	if err != nil {
		msg := fmt.Sprintf("requested article (slug=%d) not found", articleID)
		h.logger.Error().Err(err).Msg(msg)
		pp.Println(err)
		return nil, status.Error(codes.InvalidArgument, "invalid article id")
	}

	var currentUser *model.User

	// get current user if exists
	userID, err := auth.GetUserID(ctx)
	if err == nil {
		var u model.User
		err = h.db.Find(&u, userID).Error
		if err != nil {
			msg := fmt.Sprintf("token is valid but the user not found")
			h.logger.Error().Err(err).Msg(msg)
			return nil, status.Error(codes.NotFound, msg)
		}
		currentUser = &u
	}

	// bind article
	var pa pb.Article
	err = a.BindTo(&pa, currentUser, h.db)
	if err != nil {
		msg := "failed to convert model.User to pb.User"
		h.logger.Error().Err(err).Msg(msg)
		return nil, status.Error(codes.Aborted, "internal server error")
	}

	return &pb.ArticleResponse{Article: &pa}, nil
}

// GetArticles gets recent articles globally
func (h *Handler) GetArticles(ctx context.Context, req *pb.GetArticlesRequest) (*pb.ArticlesResponse, error) {
	limitQuery := req.GetLimit()
	if limitQuery == 0 {
		limitQuery = 20
	}
	offsetQuery := req.GetOffset()

	d := h.db.Preload("Author")

	// author query (has one)
	if req.GetAuthor() != "" {
		d = d.Joins("join users on articles.user_id = users.id").
			Where("users.username = ?", req.GetAuthor())
	}

	// tag query (many to many)
	if req.GetTag() != "" {
		d = d.Joins(
			"join article_tags on articles.id = article_tags.article_id "+
				"join tags on tags.id = article_tags.tag_id").
			Where("tags.name = ?", req.GetTag())
	}

	// TODO: favorite query

	// offset query, limit query
	d = d.Offset(offsetQuery).Limit(limitQuery)

	var as []model.Article
	if err := d.Find(&as).Error; err != nil {
		h.logger.Error().Err(err).Msg("failed to search articles in the database")
		pp.Println(err)
		return nil, status.Error(codes.Aborted, "internal server error")
	}

	pas := make([]*pb.Article, 0, len(as))
	for _, a := range as {
		var pa pb.Article
		err := a.BindTo(&pa, nil, nil)
		if err != nil {
			h.logger.Error().Err(err).Msg("failed to bind model.Article to pb.Article")
			return nil, status.Error(codes.Aborted, "internal server error")
		}
		pas = append(pas, &pa)
	}

	return &pb.ArticlesResponse{Articles: pas}, nil
}
