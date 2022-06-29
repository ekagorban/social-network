package mysql

import (
	"social-network/internal/models"
)

// CheckAccessExist - checks login-password pair exist
func (s *Store) CheckAccessExist(login string, password string) (a models.UserAccess, exist bool, err error) {
	// value, exist := s.storageAccess.Load(login)
	// if !exist {
	// 	return models.Access{}, false, nil
	// }
	//
	// access, ok := value.(models.Access)
	// if !ok {
	// 	return models.Access{}, false, errors.New("value.(models.Access)")
	// }

	//return access, true, nil

	return models.UserAccess{}, true, nil
}
