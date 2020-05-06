package model

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)

// Comment model
type Comment struct {
	gorm.Model
	Body      string `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	Author    User   `gorm:"foreignkey:UserID"`
	ArticleID uint   `gorm:"not null"`
	Article   Article
}

// Validate validates fields of comment model
func (c Comment) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Body,
			validation.Required,
		),
	)
}

// ProtoComment generates proto comment model from article
func (c *Comment) ProtoComment() *pb.Comment {
	pc := pb.Comment{
		Id:   fmt.Sprintf("%d", c.ID),
		Body: c.Body,
	}

	// article dates
	pc.CreatedAt, _ = ptypes.TimestampProto(c.CreatedAt)
	pc.UpdatedAt, _ = ptypes.TimestampProto(c.UpdatedAt)

	return &pc
}
