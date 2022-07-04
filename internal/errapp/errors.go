package errapp

import "errors"

var (
	UserDataNotFound   = errors.New("user data not found")
	AccessDataNotFound = errors.New("access data not found")
	LoginExist         = errors.New("login exist")
	InvalidToken       = errors.New("invalid token")
	EmptyToken         = errors.New("empty token")
)
