package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"SmartSanta/internal/delivery/events"
	"SmartSanta/internal/delivery/users"
	eventsRepo "SmartSanta/internal/repository/events"
	usersRepo "SmartSanta/internal/repository/users"
	eventsUC "SmartSanta/internal/usecase/events"
	usersUC "SmartSanta/internal/usecase/users"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitEvents(pool *pgxpool.Pool) *events.Handler {
	eRepo := eventsRepo.New(pool)
	uRepo := usersRepo.New(pool)
	uc := eventsUC.New(eRepo, uRepo)
	return events.New(uc)
}

func InitUsers(pool *pgxpool.Pool) *users.Handler {
	eRepo := eventsRepo.New(pool)
	uRepo := usersRepo.New(pool)
	uc := usersUC.New(uRepo, eRepo)
	return users.New(uc)
}

func InitDB() *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}

func main() {
	e := echo.New()

	corsCfg := middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil //todo fix
		},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}
	e.Use(middleware.CORSWithConfig(corsCfg))
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	dbpool := InitDB()
	defer dbpool.Close()

	eventHandlers := InitEvents(dbpool)
	userHandlers := InitUsers(dbpool)

	// Routes
	e.POST("/events", eventHandlers.CreateEvent)
	e.GET("/events/:hash", eventHandlers.GetEvent)
	e.PUT("/events", eventHandlers.StartEvent)

	e.POST("/users", userHandlers.CreateUser)
	e.GET("/users/:hash", userHandlers.GetUser)
	e.PUT("/users/preferences", userHandlers.UpdateUserPreferences)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
