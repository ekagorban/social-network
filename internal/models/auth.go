package models

import (
	"time"

	"github.com/google/uuid"
)

const UserAccessTable = `user_access`

// UserAccess - data for user sign up and sign in
type UserAccess struct {
	CreatedTime time.Time
	Login       string    // primary key, 3-20 symbols: letters, numbers
	Password    string    // 8-50 symbols: letters, numbers
	UserID      uuid.UUID // user id
}
