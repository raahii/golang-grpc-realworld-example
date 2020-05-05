package model

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User is user model
type User struct {
	gorm.Model
	Username string `gorm:"unique_index;not null"`
	Email    string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Bio      string `gorm:"not null"`
	Image    string `gorm:"not null"`
	Follows  []User `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
}

// Validate validates fields of user model
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(
			&u.Username,
			validation.Required,
			validation.Length(1, 10),
			validation.Match(regexp.MustCompile("[a-zA-Z0-9]+")),
		),
		validation.Field(
			&u.Email,
			validation.Required,
			is.Email,
		),
		validation.Field(
			&u.Password,
			validation.Required,
			validation.Length(6, 100),
		),
	)
}

// HashPassword makes password field crypted
func (u *User) HashPassword() error {
	if len(u.Password) == 0 {
		return errors.New("password should not be empty")
	}

	h, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(h)

	return nil
}

// CheckPassword checki user password correct
func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}
