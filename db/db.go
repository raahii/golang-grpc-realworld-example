package db

import (
	"io/ioutil"

	"os"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

func New() (*gorm.DB, error) {
	d, err := gorm.Open("sqlite3", "./db/data/realworld.db")
	if err != nil {
		return nil, err
	}

	d.DB().SetMaxIdleConns(3)
	d.LogMode(false)

	return d, nil
}

func NewTestDB() (*gorm.DB, error) {
	d, err := gorm.Open("sqlite3", "../db/data/realworld_test.db")
	if err != nil {
		return nil, err
	}

	d.DB().SetMaxIdleConns(3)
	d.LogMode(false)

	return d, nil
}

func DropTestDB() error {
	if err := os.Remove("./db/data/realworld_test.db"); err != nil {
		return err
	}
	return nil
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
	)
}

func Seed(db *gorm.DB) error {
	users := struct {
		Users []model.User
	}{}

	bs, err := ioutil.ReadFile("db/seed/users.toml")
	if err != nil {
		return err
	}

	if _, err := toml.Decode(string(bs), &users); err != nil {
		return err
	}

	for _, u := range users.Users {
		if err := db.Create(&u).Error; err != nil {
			return err
		}
	}

	return nil
}
