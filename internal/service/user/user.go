package user

import (
	"database/sql"
	"errors"
	"fmt"

	"social-network/internal/errapp"
	"social-network/internal/models"

	"github.com/google/uuid"
)

type Data struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name,required"`
	Surname string    `json:"surname,required"`
	Age     uint      `json:"age,required"`
	Gender  string    `json:"gender,required"`
	Hobbies string    `json:"hobbies,required"`
	City    string    `json:"city,required"`
}

type Service interface {
	GetAll() (users []Data, err error)
	GetOne(id uuid.UUID) (user Data, err error)
	UpdateOne(id uuid.UUID, user Data) (err error)
	GetUserFriends(id uuid.UUID) (users []Data, err error)
	AddUserFriend(userID uuid.UUID, frienID uuid.UUID) (err error)
}

type Storage interface {
	Users() (users []models.UserData, err error)
	User(id uuid.UUID) (user models.UserData, err error)
	UpdateUser(id uuid.UUID, user models.UserData) (err error)
	Friends(id uuid.UUID) (users []models.UserData, err error)
	AddFriend(userID uuid.UUID, friendID uuid.UUID) (err error)
}

type service struct {
	storage Storage
}

func NewService(store Storage) Service {
	return &service{store}
}

func (s *service) GetAll() (users []Data, err error) {
	usersModel, err := s.storage.Users()
	if err != nil {
		return nil, fmt.Errorf("s.storage.Users error: %v", err)
	}

	users = make([]Data, len(usersModel))
	for i, userModel := range usersModel {
		users[i] = Data{
			ID:      userModel.ID,
			Name:    userModel.Name,
			Surname: userModel.Surname,
			Age:     userModel.Age,
			Gender:  userModel.Gender,
			Hobbies: userModel.Hobbies,
			City:    userModel.City,
		}
	}
	return users, nil
}

func (s *service) GetOne(id uuid.UUID) (user Data, err error) {
	userModel, err := s.storage.User(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Data{}, errapp.UserDataNotFound
		}

		return Data{}, fmt.Errorf("s.storage.User error: %v", err)
	}

	user = Data{
		ID:      userModel.ID,
		Name:    userModel.Name,
		Surname: userModel.Surname,
		Age:     userModel.Age,
		Gender:  userModel.Gender,
		Hobbies: userModel.Hobbies,
		City:    userModel.City,
	}

	return user, nil
}

func (s *service) UpdateOne(id uuid.UUID, user Data) error {
	userModel := models.UserData{
		ID:      id,
		Name:    user.Name,
		Surname: user.Surname,
		Age:     user.Age,
		Gender:  user.Gender,
		Hobbies: user.Hobbies,
		City:    user.City,
	}

	err := s.storage.UpdateUser(id, userModel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errapp.UserDataNotFound
		}

		return fmt.Errorf("s.storage.UpdateUser error: %v", err)
	}

	return nil
}

func (s *service) GetUserFriends(id uuid.UUID) (users []Data, err error) {
	usersModel, err := s.storage.Friends(id)
	if err != nil {
		return nil, err
	}

	users = make([]Data, len(usersModel))
	for i, userModel := range usersModel {
		users[i] = Data{
			ID:      userModel.ID,
			Name:    userModel.Name,
			Surname: userModel.Surname,
			Age:     userModel.Age,
			Gender:  userModel.Gender,
			Hobbies: userModel.Hobbies,
			City:    userModel.City,
		}
	}
	return users, nil

}
func (s *service) AddUserFriend(userID uuid.UUID, friendID uuid.UUID) (err error) {
	return s.storage.AddFriend(userID, friendID)
}
