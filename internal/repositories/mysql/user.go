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

func (s *Store) Users() (users []models.UserData, err error) {
	query := fmt.Sprintf(`
		select 
			id,
			name,
			surname,
			age,
			gender,
			hobbies,
			city
		from %s`, models.UserDataTable)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("s.db.Query error: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("rows.Close error: %v", err)
		}
	}(rows)

	for rows.Next() {
		var user models.UserData

		err = rows.Scan(&user.ID,
			&user.Name,
			&user.Surname,
			&user.Age,
			&user.Gender,
			&user.Hobbies,
			&user.City,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan error: %v", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Store) User(id uuid.UUID) (user models.UserData, err error) {
	query := fmt.Sprintf(`
		select 
			id,
			name,
			surname,
			age,
			gender,
			hobbies,
			city
		from %s 
		where 
			id = ?`, models.UserDataTable)

	err = s.db.QueryRow(query, id).
		Scan(&user.ID,
			&user.Name,
			&user.Surname,
			&user.Age,
			&user.Gender,
			&user.Hobbies,
			&user.City)

	if err != nil {
		return models.UserData{}, err
	}

	return user, nil
}

func (s *Store) UpdateUser(id uuid.UUID, user models.UserData) (err error) {
	query := fmt.Sprintf(`
		update %s
		set
			name = ?,
			surname = ?,
			age = ?,
			gender = ?,
			hobbies = ?,
			city = ?
		where id = ?
	`, models.UserDataTable)

	err = s.db.QueryRow(query, id).
		Scan(&user.Name,
			&user.Surname,
			&user.Age,
			&user.Gender,
			&user.Hobbies,
			&user.City,
			&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Friends(id uuid.UUID) (users []models.UserData, err error) {
	query := fmt.Sprintf(`
		select 
			id,
			name,
			surname,
			age,
			gender,
			hobbies,
			city
		from %s 
		where id in (
						select friend_id 
						from %s where user_id = ?
					)
	`, models.UserDataTable, models.FriendsTable)

	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("s.db.Query error: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("rows.Close error: %v", err)
		}
	}(rows)

	for rows.Next() {
		var user models.UserData

		err = rows.Scan(&user.ID,
			&user.Name,
			&user.Surname,
			&user.Age,
			&user.Gender,
			&user.Hobbies,
			&user.City,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan error: %v", err)
		}

		users = append(users, user)
	}

	return users, nil
}
func (s *Store) AddFriend(userID uuid.UUID, friendID uuid.UUID) (err error) {
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

	log.Println("AddFriend start transaction...")

	query := fmt.Sprintf(`
		insert into %s 
			(user_id, friend_id) values
			(?, ?),
			(?, ?);`, models.FriendsTable)

	_, err = tx.ExecContext(context.Background(), query,
		userID,
		friendID,
		friendID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("s.db.ExecContext error: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx.Commit error: %v", err)
	}

	log.Println("AddFriend commit transaction success")

	return nil
}
