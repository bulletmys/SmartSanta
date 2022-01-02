package users

import (
	"SmartSanta/internal/errors"
	"SmartSanta/internal/models"
	"SmartSanta/internal/repository"
	"fmt"
	"github.com/google/uuid"
)

type UsersUC struct {
	usersRepo  repository.Users
	eventsRepo repository.Events
}

func New(usersRepo repository.Users, eventsRepo repository.Events) *UsersUC {
	return &UsersUC{usersRepo: usersRepo, eventsRepo: eventsRepo}
}

func (u *UsersUC) CreateUser(user *models.User) (string, error) {
	event, err := u.eventsRepo.GetEvent(user.EventID)
	if err != nil {
		return "", err
	}
	if event.Status != models.CREATED {
		return "", errors.WrongEventStatus
	}
	if event == nil {
		return "", errors.EventNotFound
	}

	eventUsers, err := u.usersRepo.GetUsersByEventID(event.ID)
	if err != nil {
		return "", err
	}
	if eventUsers == nil {
		user.IsAdmin = true
	}

	user.ID = uuid.New().String()
	return user.ID, u.usersRepo.CreateUser(user)
}

func (u *UsersUC) UpdateUser(user *models.User) error {
	return u.usersRepo.UpdateUser(user)
}

func (u *UsersUC) GetUser(id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid: %v", err)
	}
	user, err := u.usersRepo.GetUser(userID.String())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.UserNotFound
	}

	user.Pair, err = u.usersRepo.GetPair(user.ID, user.EventID)
	if err != nil {
		return nil, err
	}

	user.Preferences, err = u.usersRepo.GetPreferences(user.ID, user.EventID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UsersUC) UpdateUserPreferences(user *models.User) error {
	if len(user.Preferences) == 0 {
		return errors.WrongPreferencesID
	}

	userDB, err := u.usersRepo.GetUser(user.ID)
	if err != nil {
		return err
	}
	if userDB == nil {
		return errors.UserNotFound
	}

	event, err := u.eventsRepo.GetEvent(userDB.EventID)
	if err != nil {
		return err
	}
	if !(event.Status == models.STARTED || event.Status == models.FAILED) {
		return errors.WrongEventStatus
	}

	eventUsers, err := u.usersRepo.GetUsersByEventID(userDB.EventID)
	if err != nil {
		return err
	}

	if ok := compareIDS(eventUsers, user); !ok {
		return errors.WrongPreferencesID
	}

	err = u.usersRepo.UpdateUserPreferences(user)
	if err != nil {
		return err
	}

	if event.Status == models.FAILED {
		return u.eventsRepo.UpdEventStatus(event.ID, models.VOTED)
	}

	return nil
}

// Сверяем список count_id ивента из базы с предоставленным от пользователя
func compareIDS(eventIDS []models.UserWithCountID, user *models.User) bool {
	set := make(map[uint64]struct{}, len(eventIDS))

	for _, usr := range eventIDS {
		set[usr.CountID] = struct{}{}
	}

	for _, id := range user.Preferences {
		if _, ok := set[id]; !ok || id == user.CountID {
			return false
		}
	}

	return true
}
