package store

import (
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// UserStore is data access struct for user
type UserStore struct {
	db *gorm.DB
}

// NewUserStore returns a new UserStore
func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// GetByEmail finds a user from email
func (us *UserStore) GetByEmail(email string) (*model.User, error) {
	var m model.User
	if err := us.db.Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// GetByID finds a user from id
func (us *UserStore) GetByID(id uint) (*model.User, error) {
	var m model.User
	if err := us.db.Find(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// GetByUsername finds a user from username
func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	var m model.User
	if err := us.db.Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// Create create a user
func (us *UserStore) Create(m *model.User) error {
	return us.db.Create(m).Error
}

// Update update all of user fields
func (us *UserStore) Update(m *model.User) error {
	return us.db.Model(m).Update(m).Error
}

// IsFollowing returns whether user A follows user B or not
func (us *UserStore) IsFollowing(a *model.User, b *model.User) (bool, error) {
	var count int
	err := us.db.Table("follows").
		Where("from_user_id = ? AND to_user_id = ?", a.ID, b.ID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Follow create follow relashionship to User B from user A
func (us *UserStore) Follow(a *model.User, b *model.User) error {
	return us.db.Model(a).Association("Follows").Append(b).Error
}

// Unfollow delete follow relashionship to User B from user A
func (us *UserStore) Unfollow(a *model.User, b *model.User) error {
	return us.db.Model(a).Association("Follows").Delete(b).Error
}
