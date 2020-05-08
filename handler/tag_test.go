package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/raahii/golang-grpc-realworld-example/model"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)

func TestGetTags(t *testing.T) {
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

	tags := make([]string, 0, 20)
	for i := 0; i < 10; i++ {
		idStr := fmt.Sprintf("%d", i)
		a := model.Article{
			Title:       idStr,
			Description: idStr,
			Body:        idStr,
		}
		if i < 5 {
			a.Author = fooUser
		} else {
			a.Author = barUser
		}

		tag1 := uuid.New().String()
		tag2 := uuid.New().String()
		tags = append(tags, tag1)
		tags = append(tags, tag2)

		a.Tags = []model.Tag{
			model.Tag{Name: tag1},
			model.Tag{Name: tag2},
		}

		if err := h.as.Create(&a); err != nil {
			t.Fatalf("failed to create initial article record: %v", err)
		}
	}

	title := "get tags: success"
	req := &pb.Empty{}
	ctx := context.Background()
	resp, err := h.GetTags(ctx, req)

	if err != nil {
		t.Errorf("%q expected to succeed, but failed. %v", title, err)
		t.FailNow()
	}

	assert.ElementsMatch(t, resp.GetTags(), tags)
}
