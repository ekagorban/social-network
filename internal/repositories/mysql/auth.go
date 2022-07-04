package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"social-network/internal/errapp"
	"social-network/internal/models"

	"github.com/google/uuid"
)

// LoadPassword - get password by login
func (s *Store) LoadPassword(ctx context.Context, login string) (userID uuid.UUID, password string, err error) {
	query := fmt.Sprintf(`
		select user_id, password 
		from %s 
		where login = ?`, models.UserAccessTable)

	err = s.db.QueryRowContext(ctx, query, login).Scan(&userID, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, "", errapp.AccessDataNotFound
		}

		return uuid.Nil, "", fmt.Errorf("s.db.QueryRow error: %v", err)
	}

	return userID, password, nil
}
