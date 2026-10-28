package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid user",
			user: User{
				Username: "johndoe",
				Email:    "john@example.com",
				Password: "secretpassword",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			user: User{
				Username: "",
				Email:    "john@example.com",
				Password: "secretpassword",
			},
			wantErr: true,
		},
		{
			name: "invalid username characters",
			user: User{
				Username: "!!!",
				Email:    "john@example.com",
				Password: "secretpassword",
			},
			wantErr: true,
		},
		{
			name: "missing email",
			user: User{
				Username: "johndoe",
				Email:    "",
				Password: "secretpassword",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			user: User{
				Username: "johndoe",
				Email:    "not-an-email",
				Password: "secretpassword",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			user: User{
				Username: "johndoe",
				Email:    "john@example.com",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_Password(t *testing.T) {
	u := User{
		Password: "plainPassword123",
	}

	require.NoError(t, u.HashPassword())
	assert.NotEqual(t, "plainPassword123", u.Password)
	assert.True(t, u.CheckPassword("plainPassword123"))
	assert.False(t, u.CheckPassword("wrongPassword"))

	// Empty password should return error
	emptyUser := User{Password: ""}
	assert.Error(t, emptyUser.HashPassword())
}

func TestUser_ProtoConversions(t *testing.T) {
	u := User{
		Username: "alice",
		Email:    "alice@example.com",
		Bio:      "Developer",
		Image:    "https://example.com/avatar.png",
	}

	protoUser := u.ProtoUser("jwt-token-xyz")
	assert.Equal(t, u.Username, protoUser.Username)
	assert.Equal(t, u.Email, protoUser.Email)
	assert.Equal(t, "jwt-token-xyz", protoUser.Token)
	assert.Equal(t, u.Bio, protoUser.Bio)
	assert.Equal(t, u.Image, protoUser.Image)

	protoProfile := u.ProtoProfile(true)
	assert.Equal(t, u.Username, protoProfile.Username)
	assert.Equal(t, u.Bio, protoProfile.Bio)
	assert.Equal(t, u.Image, protoProfile.Image)
	assert.True(t, protoProfile.Following)
}

func TestArticle_Validate(t *testing.T) {
	tests := []struct {
		name    string
		article Article
		wantErr bool
	}{
		{
			name: "valid article",
			article: Article{
				Title: "Title",
				Body:  "Body content",
				Tags:  []Tag{{Name: "go"}},
			},
			wantErr: false,
		},
		{
			name: "missing title",
			article: Article{
				Title: "",
				Body:  "Body content",
				Tags:  []Tag{{Name: "go"}},
			},
			wantErr: true,
		},
		{
			name: "missing body",
			article: Article{
				Title: "Title",
				Body:  "",
				Tags:  []Tag{{Name: "go"}},
			},
			wantErr: true,
		},
		{
			name: "missing tags",
			article: Article{
				Title: "Title",
				Body:  "Body content",
				Tags:  nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestArticle_Overwrite(t *testing.T) {
	a := Article{
		Title:       "Original Title",
		Description: "Original Description",
		Body:        "Original Body",
	}

	a.Overwrite("New Title", "", "New Body")
	assert.Equal(t, "New Title", a.Title)
	assert.Equal(t, "Original Description", a.Description)
	assert.Equal(t, "New Body", a.Body)
}

func TestArticle_ProtoArticle(t *testing.T) {
	now := time.Now()
	a := Article{
		Title:          "Go gRPC",
		Description:    "gRPC realworld",
		Body:           "Article body",
		FavoritesCount: 5,
		Tags:           []Tag{{Name: "go"}, {Name: "grpc"}},
	}
	a.ID = 10
	a.CreatedAt = now
	a.UpdatedAt = now

	proto := a.ProtoArticle(true)
	assert.Equal(t, "10", proto.Slug)
	assert.Equal(t, "Go gRPC", proto.Title)
	assert.Equal(t, "gRPC realworld", proto.Description)
	assert.Equal(t, "Article body", proto.Body)
	assert.Equal(t, int32(5), proto.FavoritesCount)
	assert.True(t, proto.Favorited)
	assert.Equal(t, []string{"go", "grpc"}, proto.TagList)
	assert.Equal(t, now.Format(ISO8601), proto.CreatedAt)
}

func TestComment_ValidateAndProto(t *testing.T) {
	cValid := Comment{Body: "Nice article"}
	assert.NoError(t, cValid.Validate())

	cInvalid := Comment{Body: ""}
	assert.Error(t, cInvalid.Validate())

	now := time.Now()
	c := Comment{Body: "Great post!"}
	c.ID = 7
	c.CreatedAt = now
	c.UpdatedAt = now

	proto := c.ProtoComment()
	assert.Equal(t, "7", proto.Id)
	assert.Equal(t, "Great post!", proto.Body)
	assert.Equal(t, now.Format(ISO8601), proto.CreatedAt)
}
