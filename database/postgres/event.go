package postgres

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateEvent(e *models.EventRequest) error {
  var id string
	query := `INSERT INTO events (name, description, start_time, end_time, location) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return p.sql.QueryRow(query, e.Name, e.Description, e.StartTime, e.EndTime, e.Location).Scan(&id)
}

func (p *PostgresDB) GetEvent(id uuid.UUID) (*models.Event, error) {
	e := &models.Event{}
	query := `SELECT id, name, description, start_time, end_time, location FROM events WHERE id = $1`
	err := p.sql.QueryRow(query, id).Scan(&e.ID, &e.Name, &e.Description, &e.StartTime, &e.EndTime, &e.Location)
	return e, err
}

func (p *PostgresDB) UpdateEvent(e *models.Event) error {
	query := `UPDATE events SET name=$1, description=$2, start_time=$3, end_time=$4, location=$5 WHERE id=$6`
	_, err := p.sql.Exec(query, e.Name, e.Description, e.StartTime, e.EndTime, e.Location, e.ID)
	return err
}

func (p *PostgresDB) DeleteEvent(id uuid.UUID) error {
	_, err := p.sql.Exec(`DELETE FROM events WHERE id=$1`, id)
	return err
}

func (p *PostgresDB) EventExists(eventID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`
	err := p.sql.QueryRow(query, eventID).Scan(&exists)
	return exists, err
}

func (p *PostgresDB) GetAllEvents() ([]models.Event, error) {
	query := `SELECT id, name, description, start_time, end_time, location FROM events ORDER BY start_time`

	rows, err := p.sql.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.StartTime, &e.EndTime, &e.Location); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
