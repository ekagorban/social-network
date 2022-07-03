package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"social-network/internal/errapp"
	"social-network/internal/password"

	"github.com/google/uuid"
)

type SignInData struct {
	Login    string `json:"login" binding:"required,min=3,max=20"`
	Password string `json:"password,required" binding:"required,min=8,max=50"`
}

type AccessData struct {
	UserID uuid.UUID `json:"id"`
	Token  string    `json:"token,required"`
}

type Service interface {
	SignIn(data SignInData) (AccessData, error)
	CheckToken(token string) error
}

type Storage interface {
	LoadPassword(login string) (userID uuid.UUID, password string, err error)
}

type service struct {
	storage      Storage
	signingKey   []byte
	tokenExpired time.Duration
}

func NewService(store Storage) Service {
	return &service{
		storage:      store,
		signingKey:   []byte("sn-network"),
		tokenExpired: 24 * time.Hour,
	}
}

func (s *service) SignIn(data SignInData) (AccessData, error) {
	userID, passw, err := s.storage.LoadPassword(data.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AccessData{}, errapp.AccessDataNotFound
		}
		return AccessData{}, fmt.Errorf("s.storage.LoadPassword error: %v", err)
	}

	err = password.Check(passw)
	if err != nil {
		return AccessData{}, err
	}

	token, err := s.generateToken(data.Login)
	if err != nil {
		return AccessData{}, fmt.Errorf("s.generateToken error: %v", err)
	}

	accessData := AccessData{
		UserID: userID,
		Token:  token,
	}

	return accessData, nil
}

func (s *service) CheckToken(token string) error {
	err := s.parseToken(token)
	if err != nil {
		return fmt.Errorf("s.parseToken error: %v", err)
	}

	return nil
}
