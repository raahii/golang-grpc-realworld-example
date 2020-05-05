package store

import (
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// ArticleStore is data access struct for user
type ArticleStore struct {
	db *gorm.DB
}

// NewArticleStore returns a new ArticleStore
func NewArticleStore(db *gorm.DB) *ArticleStore {
	return &ArticleStore{
		db: db,
	}
}

// GetByID finds an article from id
func (s *ArticleStore) GetByID(id uint) (*model.Article, error) {
	var m model.Article
	err := s.db.Preload("Tags").Preload("Author").Find(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Create creates an article
func (s *ArticleStore) Create(m *model.Article) error {
	return s.db.Create(&m).Error
}

// Update updates an article
func (s *ArticleStore) Update(m *model.Article) error {
	return s.db.Model(&m).Update(&m).Error
}

// GetArticles get global articles
func (s *ArticleStore) GetArticles(tagName, username string, limit, offset int64) ([]model.Article, error) {
	d := s.db.Preload("Author")

	// author query (has one)
	if username != "" {
		d = d.Joins("join users on articles.user_id = users.id").
			Where("users.username = ?", username)
	}

	// tag query (many to many)
	if tagName != "" {
		d = d.Joins(
			"join article_tags on articles.id = article_tags.article_id "+
				"join tags on tags.id = article_tags.tag_id").
			Where("tags.name = ?", tagName)
	}

	// TODO: favorite query

	// offset query, limit query
	d = d.Offset(offset).Limit(limit)

	var as []model.Article
	err := d.Find(&as).Error

	return as, err
}

// GetFeedArticles returns following users' articles
func (s *ArticleStore) GetFeedArticles(userIDs []uint, limit, offset int64) ([]model.Article, error) {
	d := s.db.Preload("Author").
		Where("user_id in (?)", userIDs)

	// offset query, limit query
	d = d.Offset(offset).Limit(limit)

	var as []model.Article
	err := d.Find(&as).Error

	return as, err
}

// Delete deletes an article
func (s *ArticleStore) Delete(m *model.Article) error {
	return s.db.Delete(m).Error
}
