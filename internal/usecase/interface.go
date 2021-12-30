package usecase

import "SmartSanta/internal/models"

type Events interface {
	CreateEvent(name string) (*models.Event, error)
	GetEvent(hash string) (*models.Event, error)
	StartEvent(eventID, userID string) error
	FinishEvent(id string) error
}

type Users interface {
	CreateUser(user *models.User) (string, error)
	UpdateUser(user *models.User) error
	UpdateUserPreferences(user *models.User) error
	GetUser(id string) (*models.User, error)
}
