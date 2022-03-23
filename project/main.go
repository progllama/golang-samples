package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/jaswdr/faker"
)

type User struct {
	gorm.Model
	Name    string
	Address string
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	users := [100]User{}

	faker := faker.New()
	person := faker.Person
	for i := 0; i < 100; i++ {
		users[i] = User{
			Name:    person.Name,
			Address: person.Address,
		}
	}

	db.AutoMigrate(User{})

	db.Create(&users)
}
