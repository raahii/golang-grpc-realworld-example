package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/db"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/rs/zerolog"
)

func setUp(t *testing.T) (*Handler, func(t *testing.T)) {
	w := zerolog.ConsoleWriter{Out: ioutil.Discard}
	l := zerolog.New(w).With().Timestamp().Logger()

	d, err := db.NewTestDB()
	if err != nil {
		t.Fatal(fmt.Errorf("failed to initialize database: %w", err))
	}
	db.AutoMigrate(d)

	return New(&l, d), func(t *testing.T) {
		err := os.Remove("../db/data/realworld_test.db")
		if err != nil {
			t.Fatal(fmt.Errorf("failed to clean database: %w", err))
		}
	}
}

func TestCreateUser(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)

	tests := []struct {
		title    string
		req      *pb.CreateUserRequest
		hasError bool
	}{
		{
			"success case: foo",
			&pb.CreateUserRequest{
				User: &pb.User{
					Username: "foo",
					Email:    "foo@example.com",
					Password: "secret",
				},
			},
			false,
		},
		{
			"success case: bar",
			&pb.CreateUserRequest{
				User: &pb.User{
					Username: "bar",
					Email:    "bar@example.com",
					Password: "secret",
					Bio:      "I'm foo!",
					Image:    "https://golang.org/lib/godoc/images/go-logo-blue.svg",
				},
			},
			false,
		},
		{
			"failure case: no username",
			&pb.CreateUserRequest{
				User: &pb.User{
					Username: "",
					Email:    "foo@example.com",
					Password: "secret",
				},
			},
			true,
		},
		{
			"failure case: duplicated username",
			&pb.CreateUserRequest{
				User: &pb.User{
					Username: "foo",
					Email:    "foo@example.com",
					Password: "secret",
				},
			},
			true,
		},
		{
			"failure case: no email",
			&pb.CreateUserRequest{
				User: &pb.User{
					Username: "hoge",
					Email:    "",
					Password: "secret",
				},
			},
			true,
		},
		{
			"failure case: duplicated email",
			&pb.CreateUserRequest{
				User: &pb.User{
					Username: "hoge",
					Email:    "foo@example.com",
					Password: "secret",
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		c := context.Background()
		_, err := h.CreateUser(c, tt.req)
		if (err != nil) != tt.hasError {
			t.Errorf("%s hasError %t, but got error: %v.", tt.title, tt.hasError, err)
		}
	}
}

func TestLoginUser(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)

	fooUser := model.User{
		Username: "foo",
		Email:    "foo@example.com",
		Password: "secret",
	}
	err := fooUser.HashPassword()
	if err != nil {
		t.Fatal("Failed to hash password")
	}

	if err := h.db.Create(&fooUser).Error; err != nil {
		t.Fatalf("failed to create initial user record: %v", err)
	}

	tests := []struct {
		title    string
		req      *pb.LoginUserRequest
		expected *model.User
		hasError bool
	}{
		{
			"success case: login foo",
			&pb.LoginUserRequest{
				User: &pb.User{
					Email:    "foo@example.com",
					Password: "secret",
				},
			},
			&fooUser,
			false,
		},
	}

	for _, tt := range tests {
		c := context.Background()
		resp, err := h.LoginUser(c, tt.req)
		if (err != nil) != tt.hasError {
			t.Errorf("%q hasError %t, but got error: %v.", tt.title, tt.hasError, err)
		}

		if !tt.hasError {
			if resp.User.Username != tt.expected.Username {
				t.Errorf("%q worng Username, expected %q, got %q", tt.title, tt.expected.Username, resp.User.Username)
			}
			if resp.User.Email != tt.expected.Email {
				t.Errorf("%q worng Email, expected %q, got %q", tt.title, tt.expected.Email, resp.User.Email)
			}
			if resp.User.Bio != tt.expected.Bio {
				t.Errorf("%q worng Bio, expected %q, got %q", tt.title, tt.expected.Bio, resp.User.Bio)
			}
			if resp.User.Image != tt.expected.Image {
				t.Errorf("%q worng Image, expected %q, got %q", tt.title, tt.expected.Image, resp.User.Image)
			}
			if resp.User.Token == "" {
				t.Errorf("token must not be empety")
			}
		}
	}
}
