package signup

import (
	"errors"
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
	CreateUser(data Data) (err error)
}

type Storage interface {
	CheckLoginExist(login string) (exist bool, err error)
	CreateUser(accessData models.UserAccess, userData models.UserData) (err error)
	// CreateUserAccess(user models.UserAccess) (err error)
	// CreateUserData(user models.UserData) (err error)
}

type service struct {
	storage Storage
}

func NewService(store Storage) Service {
	return &service{store}
}

func (s *service) CreateUser(data Data) (err error) {
	err = s.checkUserLogin(data.Login)
	if err != nil {
		if errors.Is(err, errapp.LoginExist) {
			return err
		}

		return fmt.Errorf("s.checkUserLogin error: %v", err)
	}

	id := uuid.New()
	encryptedPassword := password.Encrypt(data.Password)

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

	err = s.storage.CreateUser(accessModel, dataModel)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) checkUserLogin(userLogin string) (err error) {
	exist, err := s.storage.CheckLoginExist(userLogin)
	if err != nil {
		return fmt.Errorf("s.storage.CheckLoginExist error: %v", err)
	}

	if exist {
		return fmt.Errorf("%w", errapp.LoginExist)
	}

	return nil
}
