package db

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

type Database interface {
	CreateUser(*models.User) error
	GetUser(id uuid.UUID) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(id uuid.UUID) error
	GetUsersByEvent(eventID uuid.UUID) ([]models.UserWithRole, error)

	CreateEvent(*models.Event) error
	GetEvent(id uuid.UUID) (*models.Event, error)
	UpdateEvent(*models.Event) error
	DeleteEvent(id uuid.UUID) error
	//GetEventsByAttendee(attendeeID uuid.UUID) ([]*models.Event, error)
	EventExists(eventID uuid.UUID) (bool, error)
	GetAllEvents() ([]models.Event, error)

	CreateActivity(*models.Activity) error
	GetActivity(id uuid.UUID) (*models.Activity, error)
	UpdateActivity(*models.Activity) error
	DeleteActivity(id uuid.UUID) error
	GetActivitiesByEvent(eventID uuid.UUID) ([]models.Activity, error)

	CreateAttendee(*models.Attendee) error
	GetAttendee(id uuid.UUID) (*models.Attendee, error)
	UpdateAttendee(*models.Attendee) error
	DeleteAttendee(id uuid.UUID) error
	GetAttendeesByEvent(eventID uuid.UUID) ([]models.Attendee, error)

	CreateCheckInLog(*models.CheckInLog) error
	GetCheckInLog(id uuid.UUID) (*models.CheckInLog, error)
	UpdateCheckInLog(*models.CheckInLog) error
	DeleteCheckInLog(id uuid.UUID) error
	CheckInExists(attendeeID uuid.UUID, activityID uuid.UUID) (bool, error)

	Close() error
}
