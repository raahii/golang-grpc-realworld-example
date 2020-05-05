package handler

import (
	"context"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func TestShowProfile(t *testing.T) {
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

	if err := h.us.Create(&fooUser); err != nil {
		t.Fatalf("failed to create initial user record: %v", err)
	}

	err := h.us.Follow(&fooUser, &barUser)
	if err != nil {
		t.Fatalf("failed to create initial user follow relationship: %v", err)
	}

	tests := []struct {
		title    string
		req      *pb.ShowProfileRequest
		expected *pb.Profile
		hasError bool
	}{
		{
			"show current user: success",
			&pb.ShowProfileRequest{
				Username: fooUser.Username,
			},
			&pb.Profile{
				Username:  fooUser.Username,
				Bio:       fooUser.Bio,
				Image:     fooUser.Image,
				Following: false,
			},
			false,
		},
		{
			"show following user: success",
			&pb.ShowProfileRequest{
				Username: barUser.Username,
			},
			&pb.Profile{
				Username:  barUser.Username,
				Bio:       barUser.Bio,
				Image:     barUser.Image,
				Following: true,
			},
			false,
		},
		{
			"show invalid user: invalid username",
			&pb.ShowProfileRequest{
				Username: "invalid username",
			},
			nil,
			true,
		},
	}

	token, err := auth.GenerateToken(fooUser.ID)
	if err != nil {
		t.Error(err)
	}

	for _, tt := range tests {
		ctx := ctxWithToken(context.Background(), token)

		resp, err := h.ShowProfile(ctx, tt.req)
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

		assert.Equal(t, resp.GetProfile().GetUsername(), tt.expected.Username)
		assert.Equal(t, resp.GetProfile().GetBio(), tt.expected.Bio)
		assert.Equal(t, resp.GetProfile().GetImage(), tt.expected.Image)
		assert.Equal(t, resp.GetProfile().GetFollowing(), tt.expected.GetFollowing())
	}
}

func TestFollowUser(t *testing.T) {
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
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	tests := []struct {
		title    string
		req      *pb.FollowRequest
		expected *pb.Profile
		hasError bool
	}{
		{
			"fooUser follows barUser: success",
			&pb.FollowRequest{
				Username: barUser.Username,
			},
			&pb.Profile{
				Username:  barUser.Username,
				Bio:       barUser.Bio,
				Image:     barUser.Image,
				Following: true,
			},
			false,
		},
		{
			"fooUser follows fooUser: cannnot follow myself",
			&pb.FollowRequest{
				Username: fooUser.Username,
			},
			nil,
			true,
		},
	}

	token, err := auth.GenerateToken(fooUser.ID)
	if err != nil {
		t.Error(err)
	}

	for _, tt := range tests {
		ctx := ctxWithToken(context.Background(), token)

		resp, err := h.FollowUser(ctx, tt.req)
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

		assert.Equal(t, resp.GetProfile().GetUsername(), tt.expected.GetUsername())
		assert.Equal(t, resp.GetProfile().GetBio(), tt.expected.GetBio())
		assert.Equal(t, resp.GetProfile().GetImage(), tt.expected.GetImage())
		assert.Equal(t, resp.GetProfile().GetFollowing(), tt.expected.GetFollowing())
	}
}

func TestUnfollowUser(t *testing.T) {
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
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	err := h.us.Follow(&fooUser, &barUser)
	if err != nil {
		t.Fatalf("failed to create initial user relationship: %v", err)
	}

	tests := []struct {
		title    string
		req      *pb.UnfollowRequest
		expected *pb.Profile
		hasError bool
	}{
		{
			"fooUser unfollows barUser: success",
			&pb.UnfollowRequest{
				Username: barUser.Username,
			},
			&pb.Profile{
				Username:  barUser.Username,
				Bio:       barUser.Bio,
				Image:     barUser.Image,
				Following: false,
			},
			false,
		},
		{
			"fooUser unfollows fooUser: cannnot unfollow myself",
			&pb.UnfollowRequest{
				Username: fooUser.Username,
			},
			nil,
			true,
		},
	}

	token, err := auth.GenerateToken(fooUser.ID)
	if err != nil {
		t.Error(err)
	}

	for _, tt := range tests {
		ctx := ctxWithToken(context.Background(), token)

		resp, err := h.UnfollowUser(ctx, tt.req)
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

		assert.Equal(t, resp.GetProfile().GetUsername(), tt.expected.GetUsername())
		assert.Equal(t, resp.GetProfile().GetBio(), tt.expected.GetBio())
		assert.Equal(t, resp.GetProfile().GetImage(), tt.expected.GetImage())
		assert.Equal(t, resp.GetProfile().GetFollowing(), tt.expected.GetFollowing())
	}
}
