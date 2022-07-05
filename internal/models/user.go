package models

import "github.com/google/uuid"

const UserDataTable = `user_data`
const FriendsTable = `friends`

// UserData - user information
type UserData struct {
	ID      uuid.UUID
	Name    string
	Surname string
	Age     uint
	Gender  string
	Hobbies string
	City    string
	Friends []uuid.UUID
}
