package userv1

import "time"

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
)

var SexList = []interface{}{Male, Female}

type User struct {
	ID             uint64    `json:"-"`
	UID            string    `json:"uid"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	HashedPassword *string   `json:"-"`
	Sex            Sex       `json:"sex"`
	Birthdate      time.Time `json:"birthdate"`
	Verified       bool      `json:"verified"`
	MaxSwipes      int       `json:"maxs_swipes"`
	CreatedAt      time.Time `json:"created_at"`
}
