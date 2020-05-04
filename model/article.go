package model

import "github.com/jinzhu/gorm"

// Article model
type Article struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Body        string `gorm:"not null"`
	Tags        []Tag  `gorm:"many2many:article_tags"`
	UserID      uint   `gorm:"not null"`
}

// Tag model
type Tag struct {
	gorm.Model
	Name string `gorm:"not null"`
}
