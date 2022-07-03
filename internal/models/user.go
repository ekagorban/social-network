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

// -- auto-generated definition
// create table user_data
// (
// id      varchar(36)   not null
// primary key,
// name    varchar(100)  null,
// surname varchar(100)  null,
// age     smallint      null,
// gender  char          null,
// hobbies varchar(1000) null,
// city    varchar(100)  null
// );
