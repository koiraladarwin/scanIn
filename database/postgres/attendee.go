package postgres

import "github.com/koiraladarwin/scanin/models"

func (p *PostgresDB) CreateAttendee(a models.AttendeeRequest) (models.Attendee,error) {
	query := `INSERT INTO attendees (event_id, user_id, ticket_id) VALUES ($1, $2, $3) RETURNING id, event_id, user_id, ticket_id, deleted_at;`
	var attendee models.Attendee

	err := p.sql.QueryRow(query, a.EventID, a.UserID, a.TicketID).Scan(&attendee.ID, &attendee.EventID, &attendee.UserID, &attendee.TicketID, &attendee.DeletedAt)
	if err != nil {
		return models.Attendee{},err
	}
	return attendee,nil
}
