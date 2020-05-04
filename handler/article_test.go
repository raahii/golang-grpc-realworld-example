package handler

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)

	fooUser := model.User{
		Username: "foo",
		Email:    "foo@example.com",
		Password: "secret",
	}

	for _, u := range []*model.User{&fooUser} {
		if err := h.db.Create(u).Error; err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	tests := []struct {
		title    string
		req      *pb.CreateAritcleRequest
		hasError bool
	}{
		{
			"create article: success",
			&pb.CreateAritcleRequest{
				Article: &pb.CreateAritcleRequest_Article{
					Title:       "awesome post!",
					Description: "awesome description!",
					Body:        "awesome content!",
					TagList:     []string{"foo", "bar", "piyo"},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		token, err := auth.GenerateToken(fooUser.ID)
		if err != nil {
			t.Error(err)
		}

		ctx := ctxWithToken(context.Background(), token)
		requestTime := ptypes.TimestampNow()
		resp, err := h.CreateArticle(ctx, tt.req)
		if tt.hasError {
			if err == nil {
				t.Errorf("%q expected to fail, but succeeded.", tt.title)
				t.FailNow()
			}
			continue
		}

		if !tt.hasError && err != nil {
			t.Errorf("%q expected to succeed, but failed. %v", tt.title, err)
			t.FailNow()
		}

		got := resp.GetArticle()
		expected := tt.req.GetArticle()
		assert.NotEmpty(t, got.GetSlug())
		assert.Equal(t, got.GetTitle(), expected.GetTitle())
		assert.Equal(t, got.GetDescription(), expected.GetDescription())
		assert.Equal(t, got.GetBody(), expected.GetBody())
		assert.Equal(t, got.GetTagList(), expected.GetTagList())
		assert.True(t, got.GetCreatedAt().GetNanos() > requestTime.GetNanos())
		assert.True(t, got.GetUpdatedAt().GetNanos() > requestTime.GetNanos())
		assert.False(t, got.GetFavorited())
		assert.Equal(t, got.GetFavoriteCount(), int64(0))

		author := got.GetAuthor()
		assert.Equal(t, author.GetUsername(), fooUser.Username)
		assert.Equal(t, author.GetBio(), fooUser.Bio)
		assert.Equal(t, author.GetImage(), fooUser.Image)
		assert.False(t, author.GetFollowing())
	}
}

func TestGetArticle(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)

	fooUser := model.User{
		Username: "foo",
		Email:    "foo@example.com",
		Password: "secret",
	}

	barUser := model.User{
		Username: "bar",
		Email:    "bar@example.com",
		Password: "secret",
	}

	for _, u := range []*model.User{&fooUser, &barUser} {
		if err := h.db.Create(u).Error; err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	awesomeArticle := model.Article{
		Title:       "awesome post!",
		Description: "awesome description!",
		Body:        "awesome content!",
		Tags:        []model.Tag{model.Tag{Name: "hoge"}},
		UserID:      fooUser.ID,
	}

	for _, a := range []*model.Article{&awesomeArticle} {
		if err := h.db.Create(a).Error; err != nil {
			t.Fatalf("failed to create initial article record: %v", err)
		}
	}

	tests := []struct {
		title     string
		reqUser   *model.User
		req       *pb.GetArticleRequest
		favorited bool
		following bool
		hasError  bool
	}{
		{
			"get article from unauthenticated user: success",
			nil,
			&pb.GetArticleRequest{
				Slug: string(awesomeArticle.ID),
			},
			false,
			false,
			false,
		},
		// {
		// 	"get article from barUser: success",
		// 	nil,
		// 	&pb.GetArticleRequest{
		// 		Slug: string(awesomeArticle.ID),
		// 	},
		// 	&pb.Article{
		// 		Title:       awesomeArticle.Title,
		// 		Description: awesomeArticle.Description,
		// 		Body:        awesomeArticle.Body,
		// 		TagList:     []string{awesomeArticle.Tags[0].Name},
		// 		Favorited:   false,
		// 		Author: &pb.Profile{
		// 			Username:  fooUser.Username,
		// 			Bio:       fooUser.Bio,
		// 			Image:     fooUser.Image,
		// 			Following: false,
		// 		},
		// 	},
		// },
	}

	for _, tt := range tests {
		ctx := context.Background()
		if tt.reqUser != nil {
			token, err := auth.GenerateToken(tt.reqUser.ID)
			if err != nil {
				t.Error(err)
			}

			ctx = ctxWithToken(ctx, token)
		}

		resp, err := h.GetArticle(ctx, tt.req)
		if tt.hasError {
			if err == nil {
				t.Errorf("%q expected to fail, but succeeded.", tt.title)
				t.FailNow()
			}
			continue
		}

		if !tt.hasError && err != nil {
			t.Errorf("%q expected to succeed, but failed. %v", tt.title, err)
			t.FailNow()
		}

		got := resp.GetArticle()
		assert.Equal(t, string(awesomeArticle.ID), got.GetSlug())
		assert.Equal(t, awesomeArticle.Title, got.GetTitle())
		assert.Equal(t, awesomeArticle.Description, got.GetDescription())
		assert.Equal(t, awesomeArticle.Body, got.GetBody())

		tags := make([]string, 0, len(awesomeArticle.Tags))
		for _, t := range awesomeArticle.Tags {
			tags = append(tags, t.Name)
		}
		assert.ElementsMatch(t, tags, got.GetTagList())
		assert.Equal(t, tt.favorited, got.GetFavorited())
		assert.Equal(t, int64(0), got.GetFavoriteCount())

		author := got.GetAuthor()
		assert.Equal(t, fooUser.Username, author.GetUsername())
		assert.Equal(t, fooUser.Bio, author.GetBio())
		assert.Equal(t, fooUser.Image, author.GetImage())
		assert.Equal(t, tt.following, author.GetFollowing())
	}
}
