package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)

	// create user
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
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	// create article
	awesomeArticle := model.Article{
		Title:       "awesome post!",
		Description: "awesome description!",
		Body:        "awesome content!",
		Tags:        []model.Tag{model.Tag{Name: "hoge"}},
		Author:      fooUser,
	}

	for _, a := range []*model.Article{&awesomeArticle} {
		if err := h.as.Create(a); err != nil {
			t.Fatalf("failed to create initial article record: %v", err)
		}
	}

	tests := []struct {
		title    string
		reqUser  *model.User
		req      *pb.CreateCommentRequest
		hasError bool
	}{
		{
			"create comment to awesome article: success",
			&barUser,
			&pb.CreateCommentRequest{
				Slug: fmt.Sprintf("%d", awesomeArticle.ID),
				Comment: &pb.CreateCommentRequest_Comment{
					Body: "Nice article! It helped me a lot!",
				},
			},
			false,
		},
		{
			"create comment from unauthenticated user: failed",
			nil,
			nil,
			true,
		},
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

		requestTime := ptypes.TimestampNow()
		resp, err := h.CreateComment(ctx, tt.req)
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

		got := resp.GetComment()
		assert.NotEmpty(t, got.GetId())
		assert.True(t, got.GetCreatedAt().GetNanos() > requestTime.GetNanos())
		assert.True(t, got.GetUpdatedAt().GetNanos() > requestTime.GetNanos())
		assert.Equal(t, got.GetBody(), tt.req.GetComment().GetBody())

		author := got.GetAuthor()
		assert.Equal(t, tt.reqUser.Username, author.GetUsername())
		assert.Equal(t, tt.reqUser.Bio, author.GetBio())
		assert.Equal(t, tt.reqUser.Image, author.GetImage())
		assert.False(t, author.GetFollowing())
	}
}

func TestGetComments(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)

	// create user
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

	piyoUser := model.User{
		Username: "piyo",
		Email:    "piyo@example.com",
		Password: "secret",
	}

	for _, u := range []*model.User{&fooUser, &barUser, &piyoUser} {
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	// create article
	awesomeArticle := model.Article{
		Title:       "awesome post!",
		Description: "awesome description!",
		Body:        "awesome content!",
		Tags:        []model.Tag{model.Tag{Name: "hoge"}},
		Author:      fooUser,
	}

	for _, a := range []*model.Article{&awesomeArticle} {
		if err := h.as.Create(a); err != nil {
			t.Fatalf("failed to create initial article record: %v", err)
		}
	}

	// comment articles
	comments := make([]*model.Comment, 0, 10)
	for i := 0; i < 10; i++ {
		var u model.User
		var b string
		if i%2 == 0 {
			u = barUser
			b = "ping"
		} else {
			u = piyoUser
			b = "pong"
		}
		c := model.Comment{
			Body:   b,
			Author: u,
		}
		if err := h.as.CreateComment(&c); err != nil {
			t.Fatalf("failed to create initial article comments: %v", err)
		}

		comments = append(comments, &c)
	}

	tests := []struct {
		title    string
		reqUser  *model.User
		req      *pb.GetCommentsRequest
		hasError bool
	}{
		{
			"get comments of awesome article: success",
			&barUser,
			&pb.GetCommentsRequest{
				Slug: fmt.Sprintf("%d", awesomeArticle.ID),
			},
			false,
		},
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

		resp, err := h.GetComments(ctx, tt.req)
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

		assert.Len(t, resp.GetComments(), len(comments))
		for i, got := range resp.GetComments() {
			assert.Equal(t, got.GetBody(), comments[i].Body)
			assert.Equal(t, got.GetAuthor().GetUsername, comments[i].Author.Username)
		}
	}
}
