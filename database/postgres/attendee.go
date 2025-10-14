package postgres

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateAttendee(a models.AttendeeRequest) (models.Attendee, error) {

	query := `INSERT INTO attendee ( user_id, ticket_id) VALUES ($1, $2) RETURNING id,  user_id, ticket_id;`
	var attendee models.Attendee

	err := p.sql.QueryRow(query, a.UserID, a.TicketID).Scan(&attendee.ID, &attendee.UserID, &attendee.TicketID)
	if err != nil {
		return models.Attendee{}, err
	}
	return attendee, nil
}

func (p *PostgresDB) CreateAttendeeActivityEnroll(a models.AttendeeActivityRequest) (models.AttendeeActivity, error) {

	query := `INSERT INTO attendee_activity ( attendee_id, activity_id ,firebase_id) VALUES ($1, $2, $3) RETURNING id, attendee_id, activity_id;`
	var attendeeActivityLog models.AttendeeActivity

	err := p.sql.QueryRow(query, a.AttendeeID, a.ActivityID, a.FirebaseID).Scan(&attendeeActivityLog.ID, &attendeeActivityLog.AttendeeID, &attendeeActivityLog.ActivityID)
	if err != nil {
		return models.AttendeeActivity{}, err
	}
	return attendeeActivityLog, nil
}

func (p *PostgresDB) GetEventFromAttendee(attendeeID uuid.UUID) (models.Event, error) {
	query := `
  SELECT e.id, e.firebase_id, e.event_category_id, e.name, e.event_organizer, e.description, e.start_time, e.end_time, e.location 
   FROM events e
  JOIN  ticket t ON t.event_id = e.id
  JOIN attendee a ON a.ticket_id = t.id
  WHERE a.id = $1 AND a.deleted_at IS NULL;
  `

	var event models.Event
	err := p.sql.QueryRow(query, attendeeID).Scan(
		&event.ID,
		&event.FirebaseID,
		&event.EventCategoryID,
		&event.Name,
		&event.EventOrganizer,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.Location,
	)
	if err != nil {
		return models.Event{}, err
	}
	return event, nil
}

func (p *PostgresDB) GetAttendeeFromUserIdandEventId(userID, eventID uuid.UUID) (models.Attendee, error) {
	query := `SELECT a.id, a.user_id, a.ticket_id, a.deleted_at FROM attendee a
  JOIN ticket t ON a.ticket_id = t.id
  WHERE a.user_id = $1 AND t.event_id = $2 AND a.deleted_at IS NULL;`
  var attendee models.Attendee
  err := p.sql.QueryRow(query, userID, eventID).Scan(&attendee.ID, &attendee.UserID, &attendee.TicketID, &attendee.DeletedAt)
  if err != nil {
    return models.Attendee{}, err
  }
  return attendee, nil
}
