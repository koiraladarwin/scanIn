package postgres

import (
	"github.com/google/uuid"
	db "github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateAttendee(a *models.Attendee) (*models.Attendee, error) {
	query := `INSERT INTO attendees (user_id, event_id) VALUES ($1, $2) RETURNING id`

	err := p.sql.QueryRow(query, a.UserID, a.EventID).Scan(&a.ID)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil, db.ErrAlreadyExists
		}
		return nil, err
	}
	return a, nil
}

func (p *PostgresDB) GetAttendee(id uuid.UUID) (*models.Attendee, error) {
	a := &models.Attendee{}
	query := `SELECT id, user_id, event_id  FROM attendees WHERE id=$1`
	err := p.sql.QueryRow(query, id).Scan(&a.ID, &a.UserID, &a.EventID)
	return a, err
}

func (p *PostgresDB) GetAttendeesByEvent(eventID uuid.UUID) ([]models.Attendee, error) {
	var attendees []models.Attendee
	query := `SELECT id, user_id, event_id FROM attendees WHERE event_id = $1`
	rows, err := p.sql.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var attendee models.Attendee
		if err := rows.Scan(&attendee.ID, &attendee.UserID, &attendee.EventID); err != nil {
			return nil, err
		}
		attendees = append(attendees, attendee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attendees, nil
}

func (p *PostgresDB) UpdateAttendee(a *models.Attendee) error {
	query := `UPDATE attendees SET user_id=$1, event_id=$2 WHERE id=$3`
	_, err := p.sql.Exec(query, a.UserID, a.EventID, a.ID)
	return err
}

func (p *PostgresDB) DeleteAttendee(id uuid.UUID) error {
	_, err := p.sql.Exec(`DELETE FROM attendees WHERE id=$1`, id)
	return err
}


