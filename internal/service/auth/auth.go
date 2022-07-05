package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"social-network/internal/config"
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
	SignIn(ctx context.Context, data SignInData) (AccessData, error)
	CheckToken(token string) error
}

type Storage interface {
	LoadPassword(ctx context.Context, login string) (userID uuid.UUID, password string, err error)
}

type service struct {
	storage      Storage
	signingKey   []byte
	tokenExpired time.Duration
}

func NewService(appConf config.App, store Storage) Service {
	return &service{
		storage:      store,
		signingKey:   appConf.TokenSigningKey,
		tokenExpired: appConf.TokenTimeExpired,
	}
}

func (s *service) SignIn(ctx context.Context, data SignInData) (AccessData, error) {
	userID, userPassword, err := s.storage.LoadPassword(ctx, data.Login)
	if err != nil {
		if errors.Is(err, errapp.AccessDataNotFound) {
			return AccessData{}, fmt.Errorf("s.storage.LoadPassword error: %w", errapp.AccessDataNotFound)
		}
		return AccessData{}, fmt.Errorf("s.storage.LoadPassword error: %v", err)
	}

	err = password.Check(data.Password, userPassword)
	if err != nil {
		return AccessData{}, errapp.PasswordCheckError
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
		return err
	}

	return nil
}
