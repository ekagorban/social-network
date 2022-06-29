package user

import (
	"social-network/internal/errapp"
	"social-network/internal/models"

	"github.com/google/uuid"
)

// type SignInItem struct {
// 	Login    string `json:"login,required"`
// 	Password string `json:"password,required"`
// }

// type SignUpItem struct {
// 	Login    string `json:"login,required"`
// 	Password string `json:"password,required"`
//
// 	Name    string `json:"name,required"`
// 	Surname string `json:"surname,required"`
// 	Age     uint   `json:"age,required"`
// 	Gender  string `json:"gender,required"`
// 	Hobbies string `json:"hobbies,required"`
// 	City    string `json:"city,required"`
// }

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
	User(id uuid.UUID) (user models.UserData, exist bool, err error)
	//CreateUserData(user models.UserData) (err error)
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

func (s *service) GetOne(id uuid.UUID) (user Data, err error) {
	userModel, exist, err := s.storage.User(id)
	if err != nil {
		return Data{}, err
	}

	if !exist {
		return Data{}, errapp.UserDataNotFound
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

	return s.storage.UpdateUser(id, userModel)
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
