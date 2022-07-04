package signup

import (
	"context"
	"fmt"

	"social-network/internal/errapp"
	"social-network/internal/models"
	"social-network/internal/password"

	"github.com/google/uuid"
)

type Data struct {
	Login    string `json:"login" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8,max=50"`
	Name     string `json:"name" binding:"required,max=100"`
	Surname  string `json:"surname" binding:"required,max=100"`
	Age      uint   `json:"age" binding:"required,max=100"`
	Gender   string `json:"gender" binding:"required,eq=m|eq=f"`
	Hobbies  string `json:"hobbies" binding:"required,max=1000"`
	City     string `json:"city" binding:"required,max=100"`
}

type Service interface {
	CreateUser(ctx context.Context, data Data) (err error)
}

type Storage interface {
	CheckLoginExist(ctx context.Context, login string) (exist bool, err error)
	CreateUser(ctx context.Context, accessData models.UserAccess, userData models.UserData) (err error)
}

type service struct {
	storage Storage
}

func NewService(store Storage) Service {
	return &service{store}
}

func (s *service) CreateUser(ctx context.Context, data Data) error {
	exist, err := s.checkUserLoginExist(ctx, data.Login)
	if err != nil {
		return fmt.Errorf("s.checkUserLoginExist error: %v", err)
	}

	if exist {
		return errapp.LoginExist
	}

	id := uuid.New()
	encryptedPassword, err := password.Encrypt(data.Password)
	if err != nil {
		return fmt.Errorf("password.Encrypt error: %v", err)
	}

	accessModel := models.UserAccess{
		Login:    data.Login,
		Password: encryptedPassword,
		UserID:   id,
	}
	dataModel := models.UserData{
		ID:      id,
		Name:    data.Name,
		Surname: data.Surname,
		Age:     data.Age,
		Gender:  data.Gender,
		Hobbies: data.Hobbies,
		City:    data.City,
		Friends: nil,
	}

	err = s.storage.CreateUser(ctx, accessModel, dataModel)
	if err != nil {
		return fmt.Errorf("s.storage.CreateUser error: %v", err)
	}

	return nil
}

func (s *service) checkUserLoginExist(ctx context.Context, userLogin string) (bool, error) {
	exist, err := s.storage.CheckLoginExist(ctx, userLogin)
	if err != nil {
		return false, fmt.Errorf("s.storage.CheckLoginExist error: %v", err)
	}

	if exist {
		return true, nil
	}

	return false, nil
}
