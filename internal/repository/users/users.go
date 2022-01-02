package users

import (
	"SmartSanta/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *models.User) error {
	_, err := r.db.Exec(context.Background(),
		"insert into users(user_id, name, wish, is_admin, event_id) values($1, $2, $3, $4, $5)",
		user.ID,
		user.Name,
		user.Wish,
		user.IsAdmin,
		user.EventID,
	)

	return err
}

func (r *Repository) UpdateUser(user *models.User) error {
	_, err := r.db.Exec(context.Background(),
		"update users set name = $1, wish = $2 where user_id::text = $3",
		user.Name,
		user.Wish,
		user.ID,
	)

	return err
}

func (r *Repository) GetUser(id string) (*models.User, error) {
	user := models.User{}
	err := r.db.QueryRow(context.Background(),
		"select user_id, count_id, event_id, name, wish, is_admin from users where user_id::text = $1",
		id,
	).Scan(&user.ID, &user.CountID, &user.EventID, &user.Name, &user.Wish, &user.IsAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

func (r *Repository) GetPair(id, eventID string) (*models.UserShort, error) {
	user := models.UserShort{}

	query := `select u.name, u.wish
from pairs as p
         join users u
              on (p.receiver_id = u.user_id)
where (p.sender_id::text = $1
    and p.event_id::text = $2)`
	err := r.db.QueryRow(context.Background(),
		query,
		id,
		eventID,
	).Scan(&user.Name, &user.Wish)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select pair: %v", err)
	}

	return &user, nil
}

func (r *Repository) MakePairs(users []models.UserPair, eventID string) error {
	query := `insert into pairs(sender_id, receiver_id, event_id) values($1, $2, $3)`

	for _, u := range users {
		_, err := r.db.Exec(context.Background(),
			query,
			u.CountID,
			u.PairCountID,
			eventID,
		)
		if err != nil {
			return fmt.Errorf("failed to insert pair: %v", err)
		}
	}

	return nil
}

func (r *Repository) GetPreferences(id, eventID string) ([]uint64, error) {
	var pref []sql.NullInt64

	err := r.db.QueryRow(context.Background(),
		"SELECT preferences from users where user_id::text = $1 and event_id::text = $2",
		id,
		eventID,
	).Scan(pq.Array(&pref))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select preferences: %v", err)
	}

	preferences := make([]uint64, len(pref))

	for i, id := range pref {
		preferences[i] = uint64(id.Int64)
	}

	return preferences, nil
}

func (r *Repository) UpdateUserPreferences(user *models.User) error {
	_, err := r.db.Exec(context.Background(),
		"update users set preferences = $1, is_voted = true where user_id::text = $2",
		pq.Array(user.Preferences),
		user.ID,
	)

	return err
}

func (r *Repository) GetUsersByEventID(eventID string) ([]models.UserWithCountID, error) {
	query := `select name, count_id, is_admin from users where event_id::text = $1`
	rows, err := r.db.Query(context.Background(),
		query,
		eventID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select event users: %v", err)
	}

	var isAdmin bool
	users := make([]models.UserWithCountID, 0)

	defer rows.Close()

	for rows.Next() {
		var userName models.UserWithCountID

		if err := rows.Scan(
			&userName.Name,
			&userName.CountID,
			&isAdmin,
		); err != nil {
			return nil, fmt.Errorf("error while scaning event users: %v", err)
		}

		if isAdmin {
			continue
		}
		users = append(users, userName)
	}
	if len(users) == 0 && !isAdmin {
		return nil, nil
	}

	return users, nil
}

func (r *Repository) CountVoted(eventID string) (int, error) {
	var count int

	err := r.db.QueryRow(context.Background(),
		"SELECT count(name) from users where event_id::text = $1 and is_voted = true",
		eventID,
	).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to count voted: %v", err)
	}

	return count, nil
}

func (r *Repository) GetUsersWithPreferences(eventID string) ([]models.UserWithPreferences, error) {
	query := `select count_id, preferences from users where event_id::text = $1 and is_admin = false`
	rows, err := r.db.Query(context.Background(),
		query,
		eventID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to select user preferences: %v", err)
	}

	users := make([]models.UserWithPreferences, 0)
	defer rows.Close()

	for rows.Next() {
		var user models.UserWithPreferences
		var pref []sql.NullInt64

		if err := rows.Scan(
			&user.CountID,
			pq.Array(&pref),
		); err != nil {
			return nil, fmt.Errorf("error while scaning users preferences: %v", err)
		}

		preferences := make([]uint64, len(pref))
		for i, id := range pref {
			preferences[i] = uint64(id.Int64)
		}
		user.Preferences = preferences

		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, nil
	}

	return users, nil
}
