package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"social-network/internal/errapp"
	"social-network/internal/models"

	"github.com/google/uuid"
)

func (s *Store) Users(ctx context.Context) (users []models.UserData, err error) {
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

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("s.db.QueryContext error: %v", err)
	}
	defer rowsClose(rows)

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

func (s *Store) User(ctx context.Context, id uuid.UUID) (user models.UserData, err error) {
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

	err = s.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID,
			&user.Name,
			&user.Surname,
			&user.Age,
			&user.Gender,
			&user.Hobbies,
			&user.City)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserData{}, errapp.UserDataNotFound
		}

		return models.UserData{}, fmt.Errorf("s.db.QueryRowContext error: %v", err)
	}

	return user, nil
}

func (s *Store) UpdateUser(ctx context.Context, id uuid.UUID, user models.UserData) error {
	_, err := s.User(ctx, id)
	if err != nil {
		return err
	}

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

	_, err = s.db.ExecContext(ctx, query,
		user.Name,
		user.Surname,
		user.Age,
		user.Gender,
		user.Hobbies,
		user.City,
		id)

	if err != nil {
		return fmt.Errorf("s.db.ExecContext error: %v", err)
	}

	return nil
}

func (s *Store) Friends(ctx context.Context, id uuid.UUID) (users []models.UserData, err error) {
	_, err = s.User(ctx, id)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		select
					friend.id,
					friend.name,
					friend.surname,
					friend.age,
					friend.gender,
					friend.hobbies,
					friend.city
		from %s friend
		join %s fs on friend.id = fs.friend_id
		join %s user on user.id = fs.user_id
		where user.id = ?
	`, models.UserDataTable, models.FriendsTable, models.UserDataTable)

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("s.db.QueryContext error: %v", err)
	}
	defer rowsClose(rows)

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
func (s *Store) AddFriend(ctx context.Context, userID uuid.UUID, friendID uuid.UUID) error {
	_, err := s.User(ctx, userID)
	if err != nil {
		return err
	}

	_, err = s.User(ctx, friendID)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("s.db.BeginTx error: %v", err)
	}
	defer transactionRollback(tx)

	log.Println("AddFriend start transaction...")

	query := fmt.Sprintf(`
		insert into %s 
			(user_id, friend_id) values
			(?, ?);`, models.FriendsTable)

	_, err = tx.ExecContext(ctx, query,
		userID,
		friendID,
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
