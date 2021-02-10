package models

import "math/rand"

type User struct {
	UID         int
	AccessToken string
}

var db = make(map[int]*User) //nolint:gochecknoglobals

func GetUser(id int) *User {
	return db[id]
}

func NewUser() *User {
	user := &User{UID: rand.Intn(10000)}
	db[user.UID] = user
	return user
}
