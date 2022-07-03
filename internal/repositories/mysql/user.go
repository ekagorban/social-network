package mysql

import (
	"context"
	"database/sql"
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
	// value, exist := s.storageData.Load(id)
	// if !exist {
	// 	return nil, fmt.Errorf("%w", errapp.UserDataNotFound)
	// }
	//
	// user, ok := value.(models.User)
	// if !ok {
	// 	return nil, errors.New("value.(models.User)")
	// }
	//
	// friends := make([]models.User, len(user.Friends))
	// for i, friendID := range user.Friends {
	// 	friends[i], exist, err = s.User(friendID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	//
	// return friends, nil
	return nil, nil
}
func (s *Store) AddFriend(userID uuid.UUID, friendID uuid.UUID) (err error) {
	// todo transaction and friend
	query := fmt.Sprintf(`
		insert into %s (
			user_id,
			friend_id
		) values (?, ?);`, models.FriendsTable)

	_, err = s.db.ExecContext(context.Background(), query,
		userID,
		friendID,
	)
	if err != nil {
		return fmt.Errorf("s.db.ExecContext error: %v", err)
	}

	return nil
}
