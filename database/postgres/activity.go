package postgres

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateActivity(a *models.ActivityRequest) error {
  var id string
	query := `INSERT INTO activities (event_id, name, type, start_time, end_time) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return p.sql.QueryRow(query, a.EventID, a.Name, a.Type, a.StartTime, a.EndTime).Scan(&id)
}

func (p *PostgresDB) GetActivity(id uuid.UUID) (*models.Activity, error) {
	scannedUsers := 0
	a := &models.Activity{}
	query := `SELECT id, event_id, name, type, start_time, end_time FROM activities WHERE id = $1`
	err := p.sql.QueryRow(query, id).Scan(&a.ID, &a.EventID, &a.Name, &a.Type, &a.StartTime, &a.EndTime)
	if err != nil {
		return nil, err
	}

	a.NumberOfScanedUsers = scannedUsers
	return a, err
}

func (p *PostgresDB) UpdateActivity(a *models.Activity) error {
	query := `UPDATE activities SET event_id=$1, name=$2, type=$3, start_time=$4, end_time=$5 WHERE id=$6`
	_, err := p.sql.Exec(query, a.EventID, a.Name, a.Type, a.StartTime, a.EndTime, a.ID)
	return err
}

func (p *PostgresDB) DeleteActivity(id uuid.UUID) error {
	_, err := p.sql.Exec(`DELETE FROM activities WHERE id=$1`, id)
	return err
}

func (p *PostgresDB) GetActivitiesByEvent(eventID uuid.UUID) ([]models.Activity, error) {
	activities := []models.Activity{}
	query := `SELECT id, event_id, name, type, start_time, end_time FROM activities WHERE event_id = $1`
	rows, err := p.sql.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		scannedUsers := 0

		var a models.Activity
		if err := rows.Scan(&a.ID, &a.EventID, &a.Name, &a.Type, &a.StartTime, &a.EndTime); err != nil {
			return nil, err
		}

	
query := `SELECT COUNT(*) FROM check_in_logs WHERE activity_id = $1 AND status = 'checked'`
		err = p.sql.QueryRow(query, a.ID).Scan(&scannedUsers)

		if err != nil {
			return nil, err
		}
		a.NumberOfScanedUsers = scannedUsers

		activities = append(activities, a)
	}

	return activities, nil
}
