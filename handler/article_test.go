package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/raahii/golang-grpc-realworld-example/auth"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func dateStringToUnix(d string) (int64, error) {
	t, err := time.Parse(model.ISO8601, d)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

func TestCreateArticle(t *testing.T) {
	h, cleaner := setUp(t)
	defer cleaner(t)

	fooUser := model.User{
		Username: "foo",
		Email:    "foo@example.com",
		Password: "secret",
	}

	for _, u := range []*model.User{&fooUser} {
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	tests := []struct {
		title    string
		reqUser  *model.User
		req      *pb.CreateAritcleRequest
		hasError bool
	}{
		{
			"create article: success",
			&fooUser,
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

	requestTime := time.Now().Unix() - 1
	for _, tt := range tests {
		ctx := context.Background()
		if tt.reqUser != nil {
			token, err := auth.GenerateToken(tt.reqUser.ID)
			if err != nil {
				t.Error(err)
			}

			ctx = ctxWithToken(ctx, token)
		}

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
		assert.Equal(t, expected.GetTitle(), got.GetTitle())
		assert.Equal(t, expected.GetDescription(), got.GetDescription())
		assert.Equal(t, expected.GetBody(), got.GetBody())
		assert.Equal(t, expected.GetTagList(), got.GetTagList())
		assert.True(t, got.GetFavorited())
		assert.Equal(t, int32(0), got.GetFavoritesCount())

		ct, err := dateStringToUnix(got.GetCreatedAt())
		if err != nil {
			t.Error(err)
		}
		ut, err := dateStringToUnix(got.GetUpdatedAt())
		if err != nil {
			t.Error(err)
		}
		assert.True(t, ct > requestTime)
		assert.True(t, ut > requestTime)

		author := got.GetAuthor()
		assert.Equal(t, fooUser.Username, author.GetUsername())
		assert.Equal(t, fooUser.Bio, author.GetBio())
		assert.Equal(t, fooUser.Image, author.GetImage())
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

	for _, u := range []*model.User{&barUser, &fooUser} {
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	err := h.us.Follow(&barUser, &fooUser)
	if err != nil {
		t.Fatalf("failed to create initial user relationship: %v", err)
	}

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

	if err := h.as.AddFavorite(&awesomeArticle, &barUser); err != nil {
		t.Fatalf("failed to create initial favorite articles: %v", err)
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
				Slug: fmt.Sprintf("%d", awesomeArticle.ID),
			},
			false,
			false,
			false,
		},
		{
			"get article from barUser: success",
			&barUser,
			&pb.GetArticleRequest{
				Slug: fmt.Sprintf("%d", awesomeArticle.ID),
			},
			true,
			true,
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
		assert.Equal(t, fmt.Sprintf("%d", awesomeArticle.ID), got.GetSlug())
		assert.Equal(t, awesomeArticle.Title, got.GetTitle())
		assert.Equal(t, awesomeArticle.Description, got.GetDescription())
		assert.Equal(t, awesomeArticle.Body, got.GetBody())

		tags := make([]string, 0, len(awesomeArticle.Tags))
		for _, t := range awesomeArticle.Tags {
			tags = append(tags, t.Name)
		}
		assert.ElementsMatch(t, tags, got.GetTagList())
		assert.Equal(t, tt.favorited, got.GetFavorited())
		assert.Equal(t, int32(1), got.GetFavoritesCount())

		author := got.GetAuthor()
		assert.Equal(t, fooUser.Username, author.GetUsername())
		assert.Equal(t, fooUser.Bio, author.GetBio())
		assert.Equal(t, fooUser.Image, author.GetImage())
		assert.Equal(t, tt.following, author.GetFollowing())
	}
}

func TestGetArticles(t *testing.T) {
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

	reqUser := model.User{
		Username: "req",
		Email:    "req@example.com",
		Password: "secret",
	}

	for _, u := range []*model.User{&fooUser, &barUser, &reqUser} {
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	tag := model.Tag{Name: "hoge"}

	articles := make([]*model.Article, 10)
	for i := 0; i < 10; i++ {
		idStr := fmt.Sprintf("%d", i)
		a := model.Article{
			Title:       idStr,
			Description: idStr,
			Body:        idStr,
		}
		if i < 5 {
			a.Author = fooUser
			a.Tags = []model.Tag{tag}
		} else {
			a.Author = barUser
		}

		articles[10-i-1] = &a
	}

	for i, a := range articles {
		if err := h.as.Create(a); err != nil {
			t.Fatalf("failed to create initial article record: %v", err)
		}
		if i < 5 {
			if err := h.as.AddFavorite(a, &fooUser); err != nil {
				t.Fatalf("failed to create initial favorite articles: %v", err)
			}
		}
	}

	tests := []struct {
		title    string
		req      *pb.GetArticlesRequest
		expected []*model.Article
		hasError bool
	}{
		{
			"get articles with default queries",
			&pb.GetArticlesRequest{
				Tag:       "",
				Author:    "",
				Favorited: "",
				Limit:     0,
				Offset:    0,
			},
			articles,
			false,
		},
		{
			"get articles with limit and offset",
			&pb.GetArticlesRequest{
				Tag:       "",
				Author:    "",
				Favorited: "",
				Limit:     5,
				Offset:    5,
			},
			articles[5:10],
			false,
		},
		{
			"get articles with tag",
			&pb.GetArticlesRequest{
				Tag:       "hoge",
				Author:    "",
				Favorited: "",
				Limit:     0,
				Offset:    0,
			},
			articles[5:10],
			false,
		},
		{
			"get articles with author",
			&pb.GetArticlesRequest{
				Tag:       "",
				Author:    "bar",
				Favorited: "",
				Limit:     0,
				Offset:    0,
			},
			articles[0:5],
			false,
		},
		{
			"get articles with various queries",
			&pb.GetArticlesRequest{
				Tag:       "hoge",
				Author:    "foo",
				Favorited: "",
				Limit:     2,
				Offset:    1,
			},
			articles[6:8],
			false,
		},
		{
			"get articles with favorited query",
			&pb.GetArticlesRequest{
				Tag:       "",
				Author:    "",
				Favorited: "foo",
				Limit:     0,
				Offset:    0,
			},
			articles[0:5],
			false,
		},
	}

	for _, tt := range tests {
		token, err := auth.GenerateToken(reqUser.ID)
		if err != nil {
			t.Error(err)
		}

		ctx := ctxWithToken(context.Background(), token)
		resp, err := h.GetArticles(ctx, tt.req)
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

		assert.Len(t, resp.GetArticles(), len(tt.expected))
		for i := 0; i < len(tt.expected); i++ {
			got := resp.GetArticles()[i]
			expected := tt.expected[i]

			assert.Equal(t, expected.Title, got.GetTitle(), tt.title)
			assert.Equal(t, expected.Author.Username, got.GetAuthor().GetUsername(), tt.title)
		}
	}
}

func TestGetFeedArticles(t *testing.T) {
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

	reqUser := model.User{
		Username: "req",
		Email:    "req@example.com",
		Password: "secret",
	}

	for _, u := range []*model.User{&fooUser, &barUser, &reqUser} {
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	for _, u := range []*model.User{&barUser} {
		err := h.us.Follow(&reqUser, u)
		if err != nil {
			t.Fatalf("failed to create initial user relationship: %v", err)
		}
	}

	tag := model.Tag{Name: "hoge"}

	articles := make([]*model.Article, 10)
	for i := 0; i < 10; i++ {
		idStr := fmt.Sprintf("%d", i)
		a := model.Article{
			Title:       idStr,
			Description: idStr,
			Body:        idStr,
		}
		if i < 5 {
			a.Author = fooUser
			a.Tags = []model.Tag{tag}
		} else {
			a.Author = barUser
		}

		articles[10-i-1] = &a
	}

	for _, a := range articles {
		if err := h.as.Create(a); err != nil {
			t.Fatalf("failed to create initial article record: %v", err)
		}
	}

	tests := []struct {
		title    string
		reqUser  *model.User
		req      *pb.GetFeedArticlesRequest
		expected []*model.Article
		hasError bool
	}{
		{
			"get articles with default queries",
			&reqUser,
			&pb.GetFeedArticlesRequest{
				Limit:  0,
				Offset: 0,
			},
			articles[0:5],
			false,
		},
		{
			"get articles with queries",
			&reqUser,
			&pb.GetFeedArticlesRequest{
				Limit:  2,
				Offset: 1,
			},
			articles[1:3],
			false,
		},
		{
			"get articles of user has no following user",
			&fooUser,
			&pb.GetFeedArticlesRequest{
				Limit:  2,
				Offset: 1,
			},
			[]*model.Article{},
			false,
		},
	}

	for _, tt := range tests {
		token, err := auth.GenerateToken(tt.reqUser.ID)
		if err != nil {
			t.Error(err)
		}

		ctx := ctxWithToken(context.Background(), token)
		resp, err := h.GetFeedArticles(ctx, tt.req)
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

		assert.Len(t, resp.GetArticles(), len(tt.expected))
		for i := 0; i < len(resp.GetArticles()); i++ {
			got := resp.GetArticles()[i]
			expected := tt.expected[i]

			assert.Equal(t, expected.Title, got.GetTitle())
			assert.Equal(t, expected.Author.Username, got.GetAuthor().GetUsername())
		}
	}
}

func TestUpdateArticle(t *testing.T) {
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

	for _, u := range []*model.User{&fooUser} {
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	af1 := model.Article{
		Title:       "original title",
		Description: "original desc",
		Body:        "original body",
		Author:      fooUser,
		Tags:        []model.Tag{model.Tag{Name: "hoge"}},
	}

	af2 := model.Article{
		Title:       "original title",
		Description: "original desc",
		Body:        "original body",
		Author:      fooUser,
		Tags:        []model.Tag{model.Tag{Name: "hoge"}},
	}

	ab := model.Article{
		Title:       "original title",
		Description: "original desc",
		Body:        "original body",
		Author:      barUser,
		Tags:        []model.Tag{model.Tag{Name: "hoge"}},
	}

	for _, a := range []*model.Article{&af1, &af2, &ab} {
		if err := h.as.Create(a); err != nil {
			t.Fatalf("failed to create initial article record: %v", err)
		}
	}

	tests := []struct {
		title    string
		req      *pb.UpdateArticleRequest
		expected *pb.Article
		hasError bool
	}{
		{
			"update article: success",
			&pb.UpdateArticleRequest{
				Article: &pb.UpdateArticleRequest_Article{
					Slug:        fmt.Sprintf("%d", af1.ID),
					Title:       "modified title",
					Description: "modified desc",
					Body:        "modified body",
				},
			},
			&pb.Article{
				Slug:        fmt.Sprintf("%d", af1.ID),
				Title:       "modified title",
				Description: "modified desc",
				Body:        "modified body",
				Author:      fooUser.ProtoProfile(false),
				TagList:     []string{"hoge"},
			},
			false,
		},
		{
			"update article with zero-values: no changes",
			&pb.UpdateArticleRequest{
				Article: &pb.UpdateArticleRequest_Article{
					Slug:        fmt.Sprintf("%d", af2.ID),
					Title:       "",
					Description: "",
					Body:        "",
				},
			},
			&pb.Article{
				Slug:        fmt.Sprintf("%d", af2.ID),
				Title:       "original title",
				Description: "original desc",
				Body:        "original body",
				Author:      fooUser.ProtoProfile(false),
				TagList:     []string{"hoge"},
			},
			false,
		},
		{
			"update other user's article: forbidden",
			&pb.UpdateArticleRequest{
				Article: &pb.UpdateArticleRequest_Article{
					Slug:        fmt.Sprintf("%d", ab.ID),
					Title:       "modified title",
					Description: "modified desc",
					Body:        "modified body",
				},
			},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		token, err := auth.GenerateToken(fooUser.ID)
		if err != nil {
			t.Error(err)
		}

		ctx := ctxWithToken(context.Background(), token)
		resp, err := h.UpdateArticle(ctx, tt.req)
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
		assert.Equal(t, tt.expected.GetSlug(), got.GetSlug())
		assert.Equal(t, tt.expected.GetTitle(), got.GetTitle())
		assert.Equal(t, tt.expected.GetDescription(), got.GetDescription())
		assert.Equal(t, tt.expected.GetBody(), got.GetBody())
		assert.ElementsMatch(t, tt.expected.GetTagList(), got.GetTagList())

		gotAuthor := got.GetAuthor()
		expAuthor := tt.expected.GetAuthor()
		assert.Equal(t, expAuthor.GetUsername(), gotAuthor.GetUsername())
		assert.Equal(t, expAuthor.GetBio(), gotAuthor.GetBio())
		assert.Equal(t, expAuthor.GetImage(), gotAuthor.GetImage())
		assert.Equal(t, expAuthor.GetFollowing(), gotAuthor.GetFollowing())
	}
}

func TestDeleteArticle(t *testing.T) {
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

	for _, u := range []*model.User{&fooUser} {
		if err := h.us.Create(u); err != nil {
			t.Fatalf("failed to create initial user record: %v", err)
		}
	}

	af := model.Article{
		Title:       "original title",
		Description: "original desc",
		Body:        "original body",
		Author:      fooUser,
		Tags:        []model.Tag{model.Tag{Name: "hoge"}},
	}

	ab := model.Article{
		Title:       "original title",
		Description: "original desc",
		Body:        "original body",
		Author:      barUser,
		Tags:        []model.Tag{model.Tag{Name: "hoge"}},
	}

	for _, a := range []*model.Article{&af, &ab} {
		if err := h.as.Create(a); err != nil {
			t.Fatalf("failed to create initial article record: %v", err)
		}
	}

	tests := []struct {
		title    string
		req      *pb.DeleteArticleRequest
		hasError bool
	}{
		{
			"delete article: success",
			&pb.DeleteArticleRequest{
				Slug: fmt.Sprintf("%d", af.ID),
			},
			false,
		},
		{
			"delete other user's article: forbidden",
			&pb.DeleteArticleRequest{
				Slug: fmt.Sprintf("%d", ab.ID),
			},
			true,
		},
	}

	for _, tt := range tests {
		token, err := auth.GenerateToken(fooUser.ID)
		if err != nil {
			t.Error(err)
		}

		ctx := ctxWithToken(context.Background(), token)
		_, err = h.DeleteArticle(ctx, tt.req)
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
	}
}

func TestFavoriteArticle(t *testing.T) {
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

	af := model.Article{
		Title:       "original title",
		Description: "original desc",
		Body:        "original body",
		Author:      fooUser,
		Tags:        []model.Tag{model.Tag{Name: "hoge"}},
	}

	for _, a := range []*model.Article{&af} {
		if err := h.as.Create(a); err != nil {
			t.Fatalf("failed to create initial article record: %v", err)
		}
	}

	tests := []struct {
		title    string
		reqUser  *model.User
		req      *pb.FavoriteArticleRequest
		hasError bool
	}{
		{
			"favorite user's own article: success",
			&fooUser,
			&pb.FavoriteArticleRequest{
				Slug: fmt.Sprintf("%d", af.ID),
			},
			false,
		},
		{
			"favorite other user's article: success",
			&barUser,
			&pb.FavoriteArticleRequest{
				Slug: fmt.Sprintf("%d", af.ID),
			},
			false,
		},
		{
			"favorite arleady favorited article: failed",
			&barUser,
			nil,
			true,
		},
	}

	var favoritesCount int32
	for _, tt := range tests {
		token, err := auth.GenerateToken(tt.reqUser.ID)
		if err != nil {
			t.Error(err)
		}

		ctx := ctxWithToken(context.Background(), token)
		resp, err := h.FavoriteArticle(ctx, tt.req)
		if tt.hasError {
			if err == nil {
				t.Errorf("%q expected to fail, but succeeded.", tt.title)
				t.FailNow()
			}
			continue
		}
		favoritesCount++

		if !tt.hasError && err != nil {
			t.Errorf("%q expected to succeed, but failed. %v", tt.title, err)
			t.FailNow()
		}

		got := resp.GetArticle()
		assert.True(t, got.GetFavorited())
		assert.Equal(t, favoritesCount, got.GetFavoritesCount())
	}
}

func TestUnfavoriteArticle(t *testing.T) {
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

	af := model.Article{
		Title:       "original title",
		Description: "original desc",
		Body:        "original body",
		Author:      fooUser,
		Tags:        []model.Tag{model.Tag{Name: "hoge"}},
	}

	for _, a := range []*model.Article{&af} {
		if err := h.as.Create(a); err != nil {
			t.Fatalf("failed to create initial article record: %v", err)
		}
		if err := h.as.AddFavorite(&af, &fooUser); err != nil {
			t.Fatalf("failed to create initial favorite articles: %v", err)
		}
	}

	tests := []struct {
		title          string
		reqUser        *model.User
		req            *pb.UnfavoriteArticleRequest
		favoritesCount int32
		hasError       bool
	}{
		{
			"unfavorite article: success",
			&fooUser,
			&pb.UnfavoriteArticleRequest{
				Slug: fmt.Sprintf("%d", af.ID),
			},
			0,
			false,
		},
		{
			"unfavorite not favorited article: failed",
			&barUser,
			nil,
			0,
			true,
		},
	}

	for _, tt := range tests {
		token, err := auth.GenerateToken(tt.reqUser.ID)
		if err != nil {
			t.Error(err)
		}

		ctx := ctxWithToken(context.Background(), token)
		resp, err := h.UnfavoriteArticle(ctx, tt.req)
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
		assert.False(t, got.GetFavorited())
		assert.Equal(t, tt.favoritesCount, got.GetFavoritesCount())
	}
}
