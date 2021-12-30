package events

import (
	"SmartSanta/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateEvent(name, id string) (*models.Event, error) {
	event := models.Event{}
	err := r.db.QueryRow(context.Background(),
		"insert into events(event_id, name) values($1, $2) returning event_id, name, status",
		id,
		name,
	).Scan(&event.ID, &event.Name, &event.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to insert event: %v", err)
	}

	return &event, nil
}

func (r *Repository) GetEvent(id string) (*models.Event, error) {
	event := models.Event{}
	err := r.db.QueryRow(context.Background(),
		"select event_id, name, status from events where event_id::text = $1",
		id,
	).Scan(&event.ID, &event.Name, &event.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get event: %v", err)
	}

	return &event, nil
}

func (r *Repository) UpdEventStatus(id string, status models.EventStatus) error {
	_, err := r.db.Exec(context.Background(),
		"update events set status = $1 where event_id::text = $2",
		int(status),
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to upd event status: %v", err)
	}

	return nil
}
