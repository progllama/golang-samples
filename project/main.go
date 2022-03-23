package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/jaswdr/faker"
)

type User struct {
	gorm.Model
	Name    string
	Gender  string
	Address string
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	u := [100]User{}
	var users []User = u[:]

	faker := faker.New()
	person := faker.Person()
	address := faker.Address()
	for i := 0; i < 100; i++ {
		users[i] = User{
			Name:    person.Name(),
			Gender:  person.Gender(),
			Address: address.Address(),
		}
	}

	db.AutoMigrate(User{})

	// db.Create(&users)

	// db.Select("*").Find(&users)

	// for _, u := range users {
	// 	fmt.Println(u.ID, u.Name, u.Gender, u.Address)
	// }

	user := User{}
	user.ID = 10

	db.Take(&user)
	fmt.Println(user.ID, user.Name, user.Gender, user.Address)
}
