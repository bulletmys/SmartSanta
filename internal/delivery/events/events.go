package events

import (
	"SmartSanta/internal/errors"
	"SmartSanta/internal/usecase"
	"github.com/labstack/echo/v4"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

type Handler struct {
	eventUC usecase.Events
}

func New(uc usecase.Events) *Handler {
	return &Handler{uc}
}

func (h *Handler) CreateEvent(c echo.Context) error {
	eventIn := struct {
		Name string `json:"name" validate:"nonzero"`
	}{}
	if err := c.Bind(&eventIn); err != nil {
		log.Printf("failed to decode: %v", err)
		return c.JSON(http.StatusBadRequest, "failed to decode")
	}
	if err := validator.Validate(eventIn); err != nil {
		log.Printf("failed to validate: %v", err)
		return c.JSON(http.StatusBadRequest, "failed to validate")
	}
	g, err := h.eventUC.CreateEvent(eventIn.Name)
	if err != nil {
		log.Printf("failed to create event: %v", err)
		return c.JSON(http.StatusBadRequest, "failed to create event")
	}
	return c.JSON(http.StatusOK, g)
}

func (h *Handler) GetEvent(c echo.Context) error {
	hash := c.Param("hash")
	if hash == "" {
		return c.JSON(http.StatusBadRequest, "hash value is required")
	}

	event, err := h.eventUC.GetEvent(hash)
	if err == nil {
		return c.JSON(http.StatusOK, event)
	}

	log.Printf("failed to get event: %v", err)

	switch err.(type) {
	case errors.ClientError:
		return c.JSON(http.StatusBadRequest, err.Error())
	default:
		return c.JSON(http.StatusInternalServerError, "failed to get event")
	}
}

func (h *Handler) StartEvent(c echo.Context) error {
	eventStartIn := struct {
		EventID string `json:"event_id"`
		UserID  string `json:"user_id"`
	}{}

	if err := c.Bind(&eventStartIn); err != nil {
		log.Printf("failed to decode: %v", err)
		return c.JSON(http.StatusBadRequest, "failed to decode")
	}
	if err := validator.Validate(eventStartIn); err != nil {
		log.Printf("failed to validate: %v", err)
		return c.JSON(http.StatusBadRequest, "failed to validate")
	}

	err := h.eventUC.StartEvent(eventStartIn.EventID, eventStartIn.UserID)
	if err == nil {
		return c.NoContent(http.StatusOK)
	}

	log.Printf("failed to start event: %v", err)

	switch err.(type) {
	case errors.ClientError:
		return c.JSON(http.StatusBadRequest, err.Error())
	default:
		return c.JSON(http.StatusInternalServerError, "failed to start event")
	}
}
