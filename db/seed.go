package db

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	users, err := seedUsers(db)

	if err != nil {
		fmt.Println(err)
		return
	}
	seedEvents(db, users)
}

func seedUsers(db *gorm.DB) ([]*User, error) {
	users := []*User{}

	for range 10 {
		user := User{
			Name:              faker.Name(),
			Email:             faker.Email(),
			Password:          faker.Password(),
			PasswordConfirmed: false,
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			return users, err
		}
		user.Password = string(hash)
		users = append(users, &user)
	}

	user := User{
		Name:              "Test Name",
		Email:             "test@gmail.com",
		Password:          "123",
		PasswordConfirmed: false,
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return users, err
	}
	user.Password = string(hash)
	users = append(users, &user)

	db.Create(users)

	return users, nil
}

func seedEvents(db *gorm.DB, users []*User) {
	events := []*Event{}
	for range len(users) {
        random := rand.Intn(len(users) - 1)
		event := Event{
			Title:       faker.Sentence(),
			Description: faker.Sentence(),
			Place:       faker.Word(),
			StartDate:   time.Time{},
			EndDate:     time.Time{}.Add(time.Hour * 72),
			Hosts:       []*User{users[random]},
		}

		events = append(events, &event)
	}

	db.Create(events)
}
