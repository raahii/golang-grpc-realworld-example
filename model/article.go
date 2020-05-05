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
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Body        string `gorm:"not null"`
	Tags        []Tag  `gorm:"many2many:article_tags"`
	Author      User   `gorm:"foreignkey:UserID"`
	UserID      uint   `gorm:"not null"`
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

// BindTo generates pb.Article
func (a *Article) BindTo(pa *pb.Article, requestUser *User, db *gorm.DB) error {
	pa.Slug = fmt.Sprintf("%d", a.ID)
	pa.Title = a.Title
	pa.Description = a.Description
	pa.Body = a.Body
	tags := make([]string, 0, len(a.Tags))
	for _, t := range a.Tags {
		tags = append(tags, t.Name)
	}
	pa.TagList = tags

	var err error
	pa.CreatedAt, err = ptypes.TimestampProto(a.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to convert created at field: %w", err)
	}
	pa.UpdatedAt, err = ptypes.TimestampProto(a.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to convert created at field: %w", err)
	}

	pa.Favorited = false

	pa.Author = &pb.Profile{
		Username: a.Author.Username,
		Bio:      a.Author.Bio,
		Image:    a.Author.Image,
	}

	if requestUser != nil {
		var count int
		err = db.Table("follows").Where("from_user_id = ? AND to_user_id = ?", requestUser.ID, a.Author.ID).Count(&count).Error
		if err != nil {
			return err
		}
		pa.Author.Following = count >= 1
	}

	return nil
}

// ProtoArticle generates proto aritcle model from article
func (a *Article) ProtoArticle(favorited bool) *pb.Article {
	pa := pb.Article{
		Slug:        fmt.Sprintf("%d", a.ID),
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
		Favorited:   favorited,
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
