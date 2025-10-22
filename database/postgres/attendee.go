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

func (p *PostgresDB) GetEnrolledAttendee(firebaseID string) ([]models.EnrolledAttendee, error) {
	query := `
	
SELECT 
	u.auto_id AS auto_id,
	uc.tag AS attendee_category_name,
	u.full_name AS attendee_name,
	u.image_url AS attendee_image,
	e.name AS event_name,
	a.name AS session_name,
	t.name AS ticket_name
FROM attendee_activity aa
JOIN attendee at ON aa.attendee_id = at.id
JOIN users u ON u.id = at.user_id
LEFT JOIN users_category uc ON uc.id = u.users_category_id
LEFT JOIN ticket t ON t.id = at.ticket_id
LEFT JOIN activities a ON a.id = aa.activity_id
LEFT JOIN events e ON e.id = a.event_id;
	`

	rows, err := p.sql.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendees []models.EnrolledAttendee

	for rows.Next() {
		var attendee models.EnrolledAttendee
		err := rows.Scan(
			&attendee.AutoID,
			&attendee.AttendeeCategoryName,
			&attendee.AttendeeName,
			&attendee.AttendeeImage,
			&attendee.EventName,
			&attendee.SessionName,
			&attendee.TicketName,
		)
		if err != nil {
			return nil, err
		}
		attendees = append(attendees, attendee)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return attendees, nil
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

func (p *PostgresDB) 	GetAttendeeFromUserIdandTicketId(userID, ticketID uuid.UUID) (models.Attendee, error) {
	query := `SELECT id, user_id, ticket_id, deleted_at
              FROM attendee
              WHERE user_id = $1 AND ticket_id = $2 AND deleted_at IS NULL;`

	var attendee models.Attendee
	err := p.sql.QueryRow(query, userID, ticketID).Scan(
		&attendee.ID,
		&attendee.UserID,
		&attendee.TicketID,
		&attendee.DeletedAt,
	)
	if err != nil {
		return models.Attendee{}, err
	}

	return attendee, nil
}
