package models

type EventStatus int

const (
	CREATED EventStatus = iota
	STARTED
	VOTED
	FAILED
	SUCCEED
)

type Event struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Status EventStatus       `json:"status"`
	Users  []UserWithCountID `json:"users"`
	Voted  int               `json:"voted"`
}

type User struct {
	ID          string     `json:"id"`
	CountID     uint64     `json:"count_id"`
	Name        string     `json:"name"`
	Wish        string     `json:"wish"`
	IsAdmin     bool       `json:"is_admin"`
	EventID     string     `json:"event_id"`
	Preferences []uint64   `json:"preferences"`
	Pair        *UserShort `json:"pair"`
}

type UserShort struct {
	Name string `json:"name"`
	Wish string `json:"wish"`
}

type UserWithCountID struct {
	CountID uint64 `json:"count_id"`
	Name    string `json:"name"`
}

type UserWithPreferences struct {
	CountID     uint64
	Preferences []uint64
}

type UserPair struct {
	SenderID   string
	ReceiverID string
}
