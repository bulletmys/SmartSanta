package errors

import (
	"errors"
)

type ClientError error

var (
	WrongEventStatus   ClientError = errors.New("event status is incorrect for this operation")
	EventNotFound      ClientError = errors.New("event not found")
	UserNotFound       ClientError = errors.New("user not found")
	Forbidden          ClientError = errors.New("user have no rights for this")
	WrongPreferencesID ClientError = errors.New("some of preferences id is wrong")
	NoActiveUsers      ClientError = errors.New("no active users")
)
