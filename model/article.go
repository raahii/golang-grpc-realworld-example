package model

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

// Article model
type Article struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Body           string `gorm:"not null"`
	Tags           []Tag  `gorm:"many2many:article_tags"`
	Author         User   `gorm:"foreignkey:UserID"`
	UserID         uint   `gorm:"not null"`
	FavoritesCount int64  `gorm:"not null;default=0"`
	FavoritedUsers []User `gorm:"many2many:favorite_articles"`
	Comments       []Comment
}

// Validate validates fields of article model
func (a Article) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Title,
			validation.Required,
		),
		validation.Field(
			&a.Body,
			validation.Required,
		),
		validation.Field(
			&a.Tags,
			validation.Required,
		),
	)
}

// Overwrite overwrite each field if it's not zero-value
func (a *Article) Overwrite(title, description, body string) {
	if title != "" {
		a.Title = title
	}

	if description != "" {
		a.Description = description
	}

	if body != "" {
		a.Body = body
	}
}

// ProtoArticle generates proto aritcle model from article
func (a *Article) ProtoArticle(favorited bool) *pb.Article {
	pa := pb.Article{
		Slug:           fmt.Sprintf("%d", a.ID),
		Title:          a.Title,
		Description:    a.Description,
		Body:           a.Body,
		FavoritesCount: a.FavoritesCount,
		Favorited:      favorited,
	}

	// article tags
	tags := make([]string, 0, len(a.Tags))
	for _, t := range a.Tags {
		tags = append(tags, t.Name)
	}
	pa.TagList = tags

	// article dates
	pa.CreatedAt, _ = ptypes.TimestampProto(a.CreatedAt)
	pa.UpdatedAt, _ = ptypes.TimestampProto(a.UpdatedAt)

	return &pa
}

// Tag model
type Tag struct {
	gorm.Model
	Name string `gorm:"not null"`
}
