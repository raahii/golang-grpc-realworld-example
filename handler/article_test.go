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
		assert.True(t, got.GetCreatedAt().GetSeconds() > requestTime.GetSeconds())
		assert.True(t, got.GetUpdatedAt().GetSeconds() > requestTime.GetSeconds())
		assert.False(t, got.GetFavorited())
		assert.Equal(t, got.GetFavoriteCount(), 0)

		author := got.GetAuthor()
		assert.Equal(t, author.GetUsername(), fooUser.Username)
		assert.Equal(t, author.GetBio(), fooUser.Bio)
		assert.Equal(t, author.GetBio(), fooUser.Bio)
		assert.Equal(t, author.GetImage(), fooUser.Image)
		assert.False(t, author.GetFollowing())
	}
}
