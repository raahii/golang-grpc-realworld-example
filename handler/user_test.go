package handler

import (
	"context"
	"testing"
	"time"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
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
		Bio:      "I'm foo!",
		Image:    "https://golang.org/lib/godoc/images/go-logo-blue.svg",
	}

	tests := []struct {
		title    string
		req      *pb.CreateUserRequest
		expected *model.User
		hasError bool
	}{
		{
			"create fooUser: success",
			&pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "foo",
					Email:    "foo@example.com",
					Password: "secret",
				},
			},
			&fooUser,
			false,
		},
		{
			"create barUser: success",
			&pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "bar",
					Email:    "bar@example.com",
					Password: "secret",
				},
			},
			&barUser,
			false,
		},
		{
			"create fooUser: no username",
			&pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "",
					Email:    "foo@example.com",
					Password: "secret",
				},
			},
			nil,
			true,
		},
		{
			"create fooUser: username already exists",
			&pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "foo",
					Email:    "foo@example.com",
					Password: "secret",
				},
			},
			nil,
			true,
		},
		{
			"create fooUser: no email",
			&pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "foo",
					Email:    "",
					Password: "secret",
				},
			},
			nil,
			true,
		},
		{
			"create fooUser: email already exists",
			&pb.CreateUserRequest{
				User: &pb.CreateUserRequest_User{
					Username: "hoge",
					Email:    "foo@example.com",
					Password: "secret",
				},
			},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		c := context.Background()
		resp, err := h.CreateUser(c, tt.req)
		if (err != nil) != tt.hasError {
			t.Errorf("%s hasError %t, but got error: %v.", tt.title, tt.hasError, err)
			t.FailNow()
		}

		if !tt.hasError {
			if resp.User.Username != tt.expected.Username {
				t.Errorf("%q worng Username, expected %q, got %q", tt.title, tt.expected.Username, resp.User.Username)
			}
			if resp.User.Email != tt.expected.Email {
				t.Errorf("%q worng Email, expected %q, got %q", tt.title, tt.expected.Email, resp.User.Email)
			}
			if resp.User.Token == "" {
				t.Errorf("token must not be empety")
			}
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

	if err := h.us.Create(&fooUser); err != nil {
		t.Fatalf("failed to create initial user record: %v", err)
	}

	tests := []struct {
		title    string
		req      *pb.LoginUserRequest
		expected *model.User
		hasError bool
	}{
		{
			"login to fooUser: success",
			&pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "foo@example.com",
					Password: "secret",
				},
			},
			&fooUser,
			false,
		},
		{
			"login to fooUser: wrong email",
			&pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "foooo@example.com",
					Password: "secret",
				},
			},
			nil,
			true,
		},
		{
			"login to fooUser: wrong password",
			&pb.LoginUserRequest{
				User: &pb.LoginUserRequest_User{
					Email:    "foo@example.com",
					Password: "secrets",
				},
			},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		c := context.Background()
		resp, err := h.LoginUser(c, tt.req)
		if (err != nil) != tt.hasError {
			t.Errorf("%q hasError %t, but got error: %v.", tt.title, tt.hasError, err)
			t.FailNow()
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

func TestCurrentUser(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)

	fooUser := model.User{
		Username: "foo",
		Email:    "foo@example.com",
		Password: "secret",
	}

	err := fooUser.HashPassword()
	if err != nil {
		t.Fatal("failed to hash password")
	}

	if err := h.us.Create(&fooUser); err != nil {
		t.Fatalf("failed to create initial user record: %v", err)
	}

	tests := []struct {
		title    string
		now      time.Time
		expected *model.User
		hasError bool
	}{
		{
			"get fooUser: ok",
			time.Now(),
			&fooUser,
			false,
		},
		{
			"get fooUser: token expired",
			time.Unix(0, 0),
			&fooUser,
			true,
		},
	}

	for _, tt := range tests {
		token, err := auth.GenerateTokenWithTime(tt.expected.ID, tt.now)
		if err != nil {
			t.Error(err)
		}

		ctx := ctxWithToken(context.Background(), token)
		resp, err := h.CurrentUser(ctx, &pb.Empty{})
		if (err != nil) != tt.hasError {
			t.Errorf("%q hasError %t, but got error: %v.", tt.title, tt.hasError, err)
			t.FailNow()
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

func TestUpdateUser(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)

	fooUser := model.User{
		Username: "foo",
		Email:    "foo@example.com",
		Password: "secret",
	}

	for _, u := range []*model.User{&fooUser} {
		err := u.HashPassword()
		if err != nil {
			t.Fatal("failed to hash password")
		}

		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	tests := []struct {
		title    string
		user     *model.User
		req      *pb.UpdateUserRequest
		expected *model.User
		hasError bool
	}{
		{
			"update fooUser: ok",
			&fooUser,
			&pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "hoge",
					Bio:      "Hello, World!",
				},
			},
			&model.User{
				Username: "hoge",
				Email:    "foo@example.com",
				Bio:      "Hello, World!",
				Image:    "",
			},
			false,
		},
		{
			"update fooUser: ignore zero-value field",
			&fooUser,
			&pb.UpdateUserRequest{
				User: &pb.UpdateUserRequest_User{
					Username: "foo",
					Email:    "",
				},
			},
			&model.User{
				Username: "foo",
				Email:    "foo@example.com",
				Bio:      "Hello, World!",
				Image:    "",
			},
			false,
		},
	}

	for _, tt := range tests {
		token, err := auth.GenerateToken(tt.user.ID)
		if err != nil {
			t.Error(err)
		}

		ctx := ctxWithToken(context.Background(), token)
		resp, err := h.UpdateUser(ctx, tt.req)
		if (err != nil) != tt.hasError {
			t.Errorf("%q hasError %t, but got error: %v.", tt.title, tt.hasError, err)
			t.FailNow()
		}

		if !tt.hasError {
			assert.Equal(t, resp.User.Username, tt.expected.Username)
			assert.Equal(t, resp.User.Email, tt.expected.Email)
			assert.Equal(t, resp.User.Bio, tt.expected.Bio)
			assert.Equal(t, resp.User.Image, tt.expected.Image)
			assert.NotEmpty(t, resp.User.Token)
		}
	}
}
