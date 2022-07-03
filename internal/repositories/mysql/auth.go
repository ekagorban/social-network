package mysql

import (
	"fmt"

	"social-network/internal/models"

	"github.com/google/uuid"
)

// CheckAccessExist - checks login-password pair exist
func (s *Store) LoadPassword(login string) (userID uuid.UUID, password string, err error) {
	query := fmt.Sprintf(`
		select user_id, password 
		from %s 
		where 
			login = ?`, models.UserAccessTable)

	err = s.db.QueryRow(query, login).Scan(&userID, &password)
	if err != nil {
		return uuid.Nil, "", fmt.Errorf("s.db.QueryRow error: %v", err)
	}

	return userID, password, nil
}
