package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
)

func (p *PostgresDB) CreateEvent(e *models.EventRequest) error {
	var id string
	var staff_code = utils.RandomString(6)
	var admin_code = utils.RandomString(7)

	query := `INSERT INTO events (name, description, start_time, end_time, location, staff_code, admin_code) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return p.sql.QueryRow(query, e.Name, e.Description, e.StartTime, e.EndTime, e.Location, staff_code, admin_code).Scan(&id)
}

func (p *PostgresDB) GetEventsByFirebaseUser(firebaseUser string) ([]models.Event, error) {
	query := `

SELECT 
    e.id,
    e.name,
    e.description,
    e.start_time,
    e.end_time,
    e.location,
    MAX(
      CASE 
        WHEN er.isCreator = true THEN e.staff_code
        ELSE NULL
      END
    ) AS staff_code,
    CASE
      WHEN MAX(CASE WHEN er.isCreator = true OR er.canSeeAttendee THEN 1 ELSE 0 END) = 1 THEN
        COUNT(u.id)
      ELSE
        -1
    END AS number_of_participants
FROM events e
JOIN eventRoles er ON e.id = er.event_id
LEFT JOIN users u ON u.event_id = e.id
WHERE er.fireBaseId = $1
GROUP BY e.id, e.name, e.description, e.start_time, e.end_time, e.location;
`
	rows, err := p.sql.Query(query, firebaseUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Description,
			&e.StartTime,
			&e.EndTime,
			&e.Location,
			&e.StaffCode,
			&e.NumberOfParticipant, // Make sure this field exists in your models.Event
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (p *PostgresDB) GetEventByFirebaseUser(firebaseId string, eventId uuid.UUID) (*models.Event, error) {
	e := &models.Event{}

	query := `
SELECT
  e.id,
  e.name,
  e.description,
  e.start_time,
  e.end_time,
  e.location,
  CASE WHEN er.isCreator THEN e.staff_code ELSE NULL END AS staff_code,
  CASE
    WHEN er.isCreator OR er.canSeeAttendee THEN (
      SELECT COUNT(*) FROM users WHERE event_id = e.id
    )
    ELSE -1
  END AS number_of_participant
FROM events e
JOIN eventRoles er ON e.id = er.event_id
WHERE e.id = $1 AND er.fireBaseId = $2
`

	err := p.sql.QueryRow(query, eventId, firebaseId).Scan(
		&e.ID,
		&e.Name,
		&e.Description,
		&e.StartTime,
		&e.EndTime,
		&e.Location,
		&e.StaffCode,
		&e.NumberOfParticipant,
	)
	if err != nil {
		return nil, err
	}

	fmt.Printf("number: %d \n", e.NumberOfParticipant)
	return e, nil
}

func (p *PostgresDB) GetEventByStaffId(id string) (*models.Event, error) {
	e := &models.Event{}
	query := `SELECT id, name, description, start_time, end_time, location FROM events WHERE staff_code = $1`
	err := p.sql.QueryRow(query, id).Scan(&e.ID, &e.Name, &e.Description, &e.StartTime, &e.EndTime, &e.Location)
	return e, err
}

func (p *PostgresDB) GetEventByAdminId(id string) (*models.Event, error) {
	e := &models.Event{}
	query := `SELECT id, name, description, start_time, end_time, location FROM events WHERE admin_code = $1`
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

func (p *PostgresDB) GetEventIdByActivity(acitvity uuid.UUID) (uuid.UUID, error) {
  event := models.Event{}

  query := `SELECT event_id FROM activities WHERE id = $1 limit 1`
  err := p.sql.QueryRow(query, acitvity).Scan(&event.ID)
  if err != nil {
    return uuid.Nil, fmt.Errorf("failed to get event id by activity: %w", err)
  }
  if event.ID == uuid.Nil {
    return uuid.Nil, fmt.Errorf("no event found for activity id: %s", acitvity)
  }

  return event.ID, nil
}








