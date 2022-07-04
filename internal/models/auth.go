package models

import "github.com/google/uuid"

const UserAccessTable = `user_access`

// UserAccess - data for user sign up and sign in
type UserAccess struct {
	Login    string    // primary key, 3-20 symbols: letters, numbers
	Password string    // 8-50 symbols: letters, numbers
	UserID   uuid.UUID // user id
}

// -- auto-generated definition
// create table user_access
// (
// login    varchar(20) not null
// primary key,
// password varchar(50) not null,
// user_id  varchar(36) null,
// constraint user_access_login_uindex
// unique (login)
// );
