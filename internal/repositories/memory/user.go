package memory

import (
	"errors"
	"fmt"

	"social-network/internal/errapp"
	"social-network/internal/models"

	"github.com/google/uuid"
)

func (s *Store) Users() (users []models.UserData, err error) {
	var (
		success bool
		lenData int
	)
	s.storageData.Range(func(key interface{}, value interface{}) bool {
		if id, ok := key.(uuid.UUID); ok {
			if model, ok := value.(models.UserData); ok {
				user := models.UserData{
					ID:      id,
					Name:    model.Name,
					Surname: model.Surname,
					Age:     model.Age,
					Gender:  model.Gender,
					Hobbies: model.Hobbies,
					City:    model.City,
				}
				lenData++
				users = append(users, user)
				success = true
				return true
			}
		}
		success = false
		return false
	})

	if lenData == 0 {
		return nil, nil
	}

	if !success {
		return nil, errors.New("s.storageData.Range")
	}

	return users, nil
}

func (s *Store) User(id uuid.UUID) (user models.UserData, exist bool, err error) {
	value, exist := s.storageData.Load(id)
	if !exist {
		return models.UserData{}, false, nil
	}

	user, ok := value.(models.UserData)
	if !ok {
		return models.UserData{}, false, errors.New("value.(models.User)")
	}

	return user, true, nil
}

func (s *Store) CreateUserData(user models.UserData) (err error) {
	s.storageData.Store(user.ID, user)
	return nil
}

func (s *Store) UpdateUser(id uuid.UUID, user models.UserData) (err error) {
	_, exist := s.storageData.Load(id)
	if !exist {
		return fmt.Errorf("%w", errapp.UserDataNotFound)
	}

	s.storageData.Store(id, user)
	return nil
}

func (s *Store) Friends(id uuid.UUID) (users []models.UserData, err error) {
	value, exist := s.storageData.Load(id)
	if !exist {
		return nil, fmt.Errorf("%w", errapp.UserDataNotFound)
	}

	user, ok := value.(models.UserData)
	if !ok {
		return nil, errors.New("value.(models.User)")
	}

	friends := make([]models.UserData, len(user.Friends))
	for i, friendID := range user.Friends {
		friends[i], exist, err = s.User(friendID)
		if err != nil {
			return nil, err
		}
	}

	return friends, nil

}
func (s *Store) AddFriend(userID uuid.UUID, friendID uuid.UUID) (err error) {
	valueUser, exist := s.storageData.Load(userID)
	if !exist {
		return fmt.Errorf("%w", errapp.UserDataNotFound)
	}

	user, ok := valueUser.(models.UserData)
	if !ok {
		return errors.New("valueUser.(models.User)")
	}

	valueFriend, exist := s.storageData.Load(friendID)
	if !exist {
		return fmt.Errorf("%w", errapp.UserDataNotFound)
	}
	friend, ok := valueFriend.(models.UserData)
	if !ok {
		return errors.New("valueFriend.(models.User)")
	}

	user.Friends = append(user.Friends, friendID)
	friend.Friends = append(friend.Friends, userID)

	s.storageData.Store(user.ID, user)
	s.storageData.Store(friend.ID, friend)

	return nil
}
