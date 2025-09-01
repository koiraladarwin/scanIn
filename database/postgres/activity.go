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
	query := `SELECT id, event_id, name, type, start_time, end_time FROM activities WHERE id = $1 AND delete_at IS NULL`
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

func (p *PostgresDB) GetActivitiesByEvent(firebaseId string,eventID uuid.UUID) ([]models.Activity, error) {
	activities := []models.Activity{}

	query := `
SELECT
  a.id,
  a.event_id,
  a.name,
  a.type,
  a.start_time,
  a.end_time,
  CASE
   WHEN er.isCreator OR er.canSeeScanned THEN COALESCE(scanned.count, 0)
  ELSE -1
  END AS number_of_scanned_users
FROM activities a
JOIN eventRoles er ON er.event_id = a.event_id AND er.fireBaseId = $2
LEFT JOIN (
  SELECT activity_id, COUNT(*) AS count
  FROM check_in_logs
  WHERE status = 'checked'
  GROUP BY activity_id
) scanned ON scanned.activity_id = a.id
WHERE a.event_id = $1 AND a.delete_at IS NULL;
`

	rows, err := p.sql.Query(query, eventID, firebaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.ID, &a.EventID, &a.Name, &a.Type, &a.StartTime, &a.EndTime, &a.NumberOfScanedUsers); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}

	return activities, nil
}
	
