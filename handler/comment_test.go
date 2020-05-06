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
