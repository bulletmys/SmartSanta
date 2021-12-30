package users

import (
	"SmartSanta/internal/errors"
	"SmartSanta/internal/models"
	"SmartSanta/internal/usecase"
	"github.com/labstack/echo/v4"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

type Handler struct {
	userUC usecase.Users
}

type EventIn struct {
	name string
}

func New(uc usecase.Users) *Handler {
	return &Handler{uc}
}

type UserCreateIn struct {
	Name    string `json:"name"`
	Wish    string `json:"wish"`
	EventID string `json:"event_id" validate:"nonzero"`
}

type UserCreateOut struct {
	ID string `json:"id"`
}

func (u *UserCreateIn) toModel() *models.User {
	return &models.User{
		Name:    u.Name,
		Wish:    u.Wish,
		EventID: u.EventID,
	}
}

func (h *Handler) CreateUser(c echo.Context) error {
	user := &UserCreateIn{}
	if err := c.Bind(user); err != nil {
		log.Printf("failed to decode: %v", err)
		return c.JSON(http.StatusBadRequest, "failed to decode")
	}
	if err := validator.Validate(user); err != nil {
		log.Printf("failed to validate: %v", err)
		return c.JSON(http.StatusBadRequest, "failed to validate")
	}
	id, err := h.userUC.CreateUser(user.toModel())
	if err == nil {
		return c.JSON(http.StatusOK, UserCreateOut{id})
	}

	log.Printf("failed to create user: %v", err)

	switch err.(type) {
	case errors.ClientError:
		return c.JSON(http.StatusBadRequest, err.Error())
	default:
		return c.JSON(http.StatusInternalServerError, "failed to create user")
	}
}

func (h *Handler) GetUser(c echo.Context) error {
	hash := c.Param("hash")
	if hash == "" {
		return c.JSON(http.StatusBadRequest, "hash value is required")
	}

	user, err := h.userUC.GetUser(hash)
	if err == nil {
		return c.JSON(http.StatusOK, user)
	}

	log.Printf("failed to get user: %v", err)

	switch err.(type) {
	case errors.ClientError:
		return c.JSON(http.StatusBadRequest, err.Error())
	default:
		return c.JSON(http.StatusInternalServerError, "failed to get user")
	}
}

type UserPreferencesIn struct {
	UserID string   `json:"user_id" validate:"nonzero"`
	IDS    []uint64 `json:"ids" validate:"nonzero"`
}

func (u *UserPreferencesIn) toModel() *models.User {
	return &models.User{
		ID:          u.UserID,
		Preferences: u.IDS,
	}
}

func (h *Handler) UpdateUserPreferences(c echo.Context) error {
	user := &UserPreferencesIn{}
	if err := c.Bind(user); err != nil {
		log.Printf("failed to decode: %v", err)
		return c.JSON(http.StatusBadRequest, "failed to decode")
	}
	if err := validator.Validate(user); err != nil {
		log.Printf("failed to validate: %v", err)
		return c.JSON(http.StatusBadRequest, "failed to validate")
	}
	err := h.userUC.UpdateUserPreferences(user.toModel())
	if err == nil {
		return c.NoContent(http.StatusOK)
	}

	log.Printf("failed to upd user preferences: %v", err)

	switch err.(type) {
	case errors.ClientError:
		return c.JSON(http.StatusBadRequest, err.Error())
	default:
		return c.JSON(http.StatusInternalServerError, "failed to create user")
	}
}
