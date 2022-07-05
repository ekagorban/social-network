package errapp

import "errors"

var (
	UserDataNotFound   = errors.New("user data not found")
	AccessDataNotFound = errors.New("access data not found")
	PasswordCheckError = errors.New("password check error")
	LoginExist         = errors.New("login exist")
	InvalidToken       = errors.New("invalid token")
	ExpiredToken       = errors.New("expired token")
	EmptyToken         = errors.New("empty token")
)
