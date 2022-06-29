package memory

import (
	"errors"

	"social-network/internal/models"
)

// CheckLoginExist - checks login exist
func (s *Store) CheckLoginExist(login string) (exist bool, err error) {
	if _, ok := s.storageAccess.Load(login); ok {
		return true, nil
	}
	return false, nil
}

// CheckAccessExist - checks login-password pair exist
func (s *Store) CheckAccessExist(login string, password string) (a models.UserAccess, exist bool, err error) {
	value, exist := s.storageAccess.Load(login)
	if !exist {
		return models.UserAccess{}, false, nil
	}

	access, ok := value.(models.UserAccess)
	if !ok {
		return models.UserAccess{}, false, errors.New("value.(models.Access)")
	}

	return access, true, nil
}

func (s *Store) CreateUserAccess(user models.UserAccess) (err error) {
	s.storageAccess.Store(user.Login, user)
	return nil
}
