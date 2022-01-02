package events

import (
	"SmartSanta/internal/errors"
	"SmartSanta/internal/models"
	"SmartSanta/internal/repository"
	"github.com/google/uuid"
)

type EventsUC struct {
	eventsRepo repository.Events
	usersRepo  repository.Users
}

func New(eventsRepo repository.Events, usersRepo repository.Users) *EventsUC {
	return &EventsUC{eventsRepo: eventsRepo, usersRepo: usersRepo}
}

func (u *EventsUC) CreateEvent(name string) (*models.Event, error) {
	return u.eventsRepo.CreateEvent(name, uuid.New().String())
}

func (u *EventsUC) GetEvent(id string) (*models.Event, error) {
	event, err := u.eventsRepo.GetEvent(id)
	if err != nil {
		return nil, errors.EventNotFound
	}
	event.Users, err = u.usersRepo.GetUsersByEventID(id)
	if err != nil {
		return nil, err
	}

	event.Voted, err = u.usersRepo.CountVoted(id)

	return event, err
}

func (u *EventsUC) StartEvent(eventID, userID string) error {
	event, err := u.eventsRepo.GetEvent(eventID)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.EventNotFound
	}
	if event.Status != models.CREATED {
		return errors.WrongEventStatus
	}
	user, err := u.usersRepo.GetUser(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.UserNotFound
	}
	if user.IsAdmin == false {
		return errors.Forbidden
	}
	return u.eventsRepo.UpdEventStatus(eventID, models.STARTED)
}

func (u *EventsUC) FinishEvent(id string) error {
	event, err := u.eventsRepo.GetEvent(id)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.EventNotFound
	}
	if event.Status != models.STARTED {
		return errors.WrongEventStatus
	}
	return u.eventsRepo.UpdEventStatus(id, models.VOTED)
}
