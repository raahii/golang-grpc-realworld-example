package handler

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/raahii/golang-grpc-realworld-example/db"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/rs/zerolog"
)

func setUp(t *testing.T) (*Handler, func(t *testing.T)) {
	w := zerolog.ConsoleWriter{Out: os.Stderr}
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

func TestCreateuser(t *testing.T) {
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
