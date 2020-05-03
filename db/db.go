package db

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"os"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/raahii/golang-grpc-realworld-example/model"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func dsn() (string, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return "", errors.New("$DB_HOST is not set")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		return "", errors.New("$DB_USER is not set")
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return "", errors.New("$DB_PASSWORD is not set")
	}

	name := os.Getenv("DB_NAME")
	if name == "" {
		return "", errors.New("$DB_NAME is not set")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		return "", errors.New("$DB_PORT is not set")
	}

	options := "charset=utf8mb4&parseTime=True&loc=Local"

	// "user:password@host:port/dbname?option1&option2"
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?%s",
		user, password, host, port, name, options), nil
}

func New() (*gorm.DB, error) {
	s, err := dsn()
	if err != nil {
		return nil, err
	}

	var d *gorm.DB
	for i := 0; i < 10; i++ {
		d, err = gorm.Open("mysql", s)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, err
	}

	d.DB().SetMaxIdleConns(3)
	d.LogMode(false)

	return d, nil
}

func NewTestDB() (*gorm.DB, error) {
	//TODO: replace test database
	err := godotenv.Load("../env/local.env")
	if err != nil {
		return nil, err
	}
	return New()
}

func DropTestDB(d *gorm.DB) error {
	return d.DropTableIfExists("users", "follows").Error
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.User{},
	).Error
	if err != nil {
		return err
	}
	return nil
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
