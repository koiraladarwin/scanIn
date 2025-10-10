package postgres

import "github.com/koiraladarwin/scanin/models"

func (p *PostgresDB) CreateAttendee(a models.AttendeeRequest) (models.Attendee, error) {

	query := `INSERT INTO attendee ( user_id, ticket_id) VALUES ($1, $2) RETURNING id,  user_id, ticket_id;`
	var attendee models.Attendee

	err := p.sql.QueryRow(query, a.UserID, a.TicketID).Scan(&attendee.ID, &attendee.UserID, &attendee.TicketID)
	if err != nil {
		return models.Attendee{}, err
	}
	return attendee, nil
}

func (p *PostgresDB) CreateAttendeeActivityEnroll(a models.AttendeeActivity) (models.AttendeeActivity, error) {

  query := `INSERT INTO attendee_activity ( attendee_id, activity_id, ) VALUES ($1, $2) RETURNING id, attendee_id, activity_id;`
  var attendeeActivityLog models.AttendeeActivity

  err := p.sql.QueryRow(query, a.AttendeeID, a.ActivityID).Scan(&attendeeActivityLog.ID, &attendeeActivityLog.AttendeeID)
  if err != nil {
    return models.AttendeeActivity{}, err
  }
  return attendeeActivityLog, nil
}
