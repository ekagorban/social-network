package mysql

import (
	"context"
	"fmt"

	"social-network/internal/models"
)

func (s *Store) UsersByFilters(ctx context.Context, name string, surname string) (users []models.UserData, err error) {
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
		where name like ? and surname like ?`, models.UserDataTable)

	rows, err := s.db.QueryContext(ctx, query, name, surname)
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
