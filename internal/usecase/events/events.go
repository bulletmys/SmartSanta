package events

import (
	"SmartSanta/internal/algorithm"
	"SmartSanta/internal/errors"
	"SmartSanta/internal/models"
	"SmartSanta/internal/repository"
	"github.com/google/uuid"
	"log"
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

	if event.Voted == len(event.Users) && event.Status == models.VOTED {
		u.StartCountWrapper(event)
	}

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

func (u *EventsUC) StartCountWrapper(event *models.Event) {
	go u.StartCount(event)
}

func convertModelToPreferencesMap(users []models.UserWithPreferences) map[uint64][]uint64 {
	prefMap := make(map[uint64][]uint64)

	for _, u := range users {
		prefMap[u.CountID] = u.Preferences
	}

	return prefMap
}

func convertPreferencesMapToModel(prefs map[uint64]uint64) []models.UserPair {
	users := make([]models.UserPair, len(prefs))
	i := 0

	for k, v := range prefs {
		users[i] = models.UserPair{
			CountID:     k,
			PairCountID: v,
		}
		i++
	}

	return users
}

func (u *EventsUC) MarkFailed(event *models.Event) {
	err := u.eventsRepo.UpdEventStatus(event.ID, models.FAILED)
	if err != nil {
		log.Printf("failed to change event status to failed: %v", err)
	}
}

func (u *EventsUC) StartCount(event *models.Event) {
	err := u.FinishEvent(event.ID)
	if err != nil {
		log.Printf("failed to finish event: %v", err)
		u.MarkFailed(event)
		return
	}
	users, err := u.usersRepo.GetUsersWithPreferences(event.ID)
	if err != nil {
		log.Printf("failed to get users preferences: %v", err)
		u.MarkFailed(event)
		return
	}

	pairs := algorithm.CountPreferences(convertModelToPreferencesMap(users))
	if pairs == nil {
		log.Printf("failed to count pairs: %v", err)
		u.MarkFailed(event)
	}

	err = u.usersRepo.MakePairs(convertPreferencesMapToModel(pairs), event.ID)
	if err != nil {
		log.Printf("failed to make pairs: %v", err)
		u.MarkFailed(event)
		return
	}
}
