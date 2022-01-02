package repository

import (
	"SmartSanta/internal/models"
)

type Events interface {
	CreateEvent(name, id string) (*models.Event, error)
	GetEvent(id string) (*models.Event, error)
	UpdEventStatus(id string, status models.EventStatus) error
}

type Users interface {
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetUser(id string) (*models.User, error)
	GetPair(id, eventID string) (*models.UserShort, error)
	MakePairs(users []models.UserPair, eventID string) error
	GetPreferences(id, eventID string) ([]uint64, error)
	UpdateUserPreferences(user *models.User) error
	GetUsersByEventID(eventID string) ([]models.UserWithCountID, error)
	CountVoted(eventID string) (int, error)
	GetUsersWithPreferences(eventID string) ([]models.UserWithPreferences, error)
}
