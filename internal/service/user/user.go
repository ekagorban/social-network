package user

import (
	"context"
	"errors"
	"fmt"

	"social-network/internal/errapp"
	"social-network/internal/models"

	"github.com/google/uuid"
)

type Data struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name" binding:"required,max=100"`
	Surname string    `json:"surname" binding:"required,max=100"`
	Age     uint      `json:"age" binding:"required,max=100"`
	Gender  string    `json:"gender" binding:"required,eq=m|eq=f"`
	Hobbies string    `json:"hobbies" binding:"required,max=1000"`
	City    string    `json:"city" binding:"required,max=100"`
}

type Service interface {
	GetAll(ctx context.Context) (users []Data, err error)
	GetOne(ctx context.Context, id uuid.UUID) (user Data, err error)
	UpdateOne(ctx context.Context, id uuid.UUID, user Data) (err error)
	GetUserFriends(ctx context.Context, id uuid.UUID) (users []Data, err error)
	AddUserFriend(ctx context.Context, userID uuid.UUID, friendID uuid.UUID) (err error)
}

type Storage interface {
	Users(ctx context.Context) (users []models.UserData, err error)
	User(ctx context.Context, id uuid.UUID) (user models.UserData, err error)
	UpdateUser(ctx context.Context, id uuid.UUID, user models.UserData) (err error)
	Friends(ctx context.Context, id uuid.UUID) (users []models.UserData, err error)
	AddFriend(ctx context.Context, userID uuid.UUID, friendID uuid.UUID) (err error)
}

type service struct {
	storage Storage
}

func NewService(store Storage) Service {
	return &service{store}
}

func (s *service) GetAll(ctx context.Context) (users []Data, err error) {
	usersModel, err := s.storage.Users(ctx)
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

func (s *service) GetOne(ctx context.Context, id uuid.UUID) (user Data, err error) {
	userModel, err := s.storage.User(ctx, id)
	if err != nil {
		if errors.Is(err, errapp.UserDataNotFound) {
			return Data{}, fmt.Errorf("s.storage.User error: %w", errapp.UserDataNotFound)
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

func (s *service) UpdateOne(ctx context.Context, id uuid.UUID, user Data) error {
	userModel := models.UserData{
		ID:      id,
		Name:    user.Name,
		Surname: user.Surname,
		Age:     user.Age,
		Gender:  user.Gender,
		Hobbies: user.Hobbies,
		City:    user.City,
	}

	err := s.storage.UpdateUser(ctx, id, userModel)
	if err != nil {
		if errors.Is(err, errapp.UserDataNotFound) {
			return fmt.Errorf("s.storage.UpdateUser error: %w", errapp.UserDataNotFound)
		}

		return fmt.Errorf("s.storage.UpdateUser error: %v", err)
	}

	return nil
}

func (s *service) GetUserFriends(ctx context.Context, id uuid.UUID) (users []Data, err error) {
	usersModel, err := s.storage.Friends(ctx, id)
	if err != nil {
		if errors.Is(err, errapp.UserDataNotFound) {
			return nil, fmt.Errorf("s.storage.Friends error: %w", errapp.UserDataNotFound)
		}
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
func (s *service) AddUserFriend(ctx context.Context, userID uuid.UUID, friendID uuid.UUID) error {
	err := s.storage.AddFriend(ctx, userID, friendID)
	if err != nil {
		if errors.Is(err, errapp.UserDataNotFound) {
			return fmt.Errorf("s.storage.AddFriend error: %w", errapp.UserDataNotFound)
		}

		return fmt.Errorf("s.storage.UpdateUser error: %v", err)
	}

	return nil
}
