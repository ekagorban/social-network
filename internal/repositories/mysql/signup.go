package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"social-network/internal/models"

	"github.com/google/uuid"
)

// CheckLoginExist - checks login exist
func (s *Store) CheckLoginExist(login string) (exist bool, err error) {
	query := fmt.Sprintf(`
		select user_id 
		from %s 
		where login = ?`, models.UserAccessTable)

	var userID uuid.UUID

	err = s.db.QueryRow(query, login).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *Store) CreateUser(access models.UserAccess, user models.UserData) (err error) {
	ctx := context.Background()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("s.db.BeginTx error: %v", err)
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				log.Printf("tx.Rollback error: %v", err)
			}
		}
	}(tx)

	log.Println("CreateUser start transaction...")

	err = s.createUserData(ctx, tx, user)
	if err != nil {
		return fmt.Errorf("s.createUserData error: %v", err)
	}

	err = s.createAccessData(ctx, tx, access)
	if err != nil {
		return fmt.Errorf("s.createAccessData error: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit error: %v", err)
	}

	log.Println("CreateUser commit transaction success")

	return nil
}

func (s *Store) createAccessData(ctx context.Context, tx *sql.Tx, access models.UserAccess) (err error) {
	query := fmt.Sprintf(`
		insert into %s (
			login,
			password,
			user_id
		) values (?, ?, ?)`, models.UserAccessTable)

	_, err = tx.ExecContext(ctx, query,
		access.Login,
		access.Password,
		access.UserID,
	)
	if err != nil {
		return fmt.Errorf("tx.ExecContext error: %v", err)
	}

	return nil
}

func (s *Store) createUserData(ctx context.Context, tx *sql.Tx, user models.UserData) (err error) {
	query := fmt.Sprintf(`
		insert into %s (
			id, 
			name, 
			surname, 
			age, 
			gender, 
			hobbies, 
			city
		) values (?, ?, ?, ?, ?, ?, ?)`, models.UserDataTable)

	_, err = tx.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Surname,
		user.Age,
		user.Gender,
		user.Hobbies,
		user.City,
	)
	if err != nil {
		return fmt.Errorf("tx.ExecContext error: %v", err)
	}

	return nil
}
