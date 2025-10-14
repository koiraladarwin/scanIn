package postgres

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateActivity(a *models.ActivityCreateRequest) (*models.ActivityWithScannedUser, error) {
	var activity models.ActivityWithScannedUser
	query := `
		INSERT INTO activities (event_id, name, hall_name, start_time, end_time ,firebase_id) 
		VALUES ($1, $2, $3, $4, $5 ,$6)
		RETURNING id, event_id, name, hall_name, start_time, end_time
	`

	err := p.sql.QueryRow(
		query,
		a.EventID,
		a.Name,
		a.HallName,
		a.StartTime,
		a.EndTime,
		a.FirebaseID,
	).Scan(
		&activity.ID,
		&activity.EventID,
		&activity.Name,
		&activity.HallName,
		&activity.StartTime,
		&activity.EndTime,
	)

	if err != nil {
		return nil, err
	}

	return &activity, nil
}

func (p *PostgresDB) GetActivitiesDetails(firebaseId string) ([]models.ActivityDetails, error) {
	query := `
SELECT
    a.id,
    a.event_id,
    a.name,
    a.hall_name,
    a.start_time,
    a.end_time,
    COUNT(t.id) FILTER (WHERE t.price = 0) AS ticket_count,
    COUNT(t.id) FILTER (WHERE t.price != 0) AS invitation_count
FROM activities a
LEFT JOIN attendee_activity aa ON aa.activity_id = a.id
LEFT JOIN attendee at ON at.id = aa.attendee_id
LEFT JOIN ticket t ON t.id = at.ticket_id
WHERE a.deleted_at IS NULL
  AND a.firebase_id = $1
GROUP BY a.id;
`

	rows, err := p.sql.Query(query, firebaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.ActivityDetails

	for rows.Next() {
		var act models.ActivityDetails
		err := rows.Scan(
			&act.ID,
			&act.EventID,
			&act.Name,
			&act.HallName,
			&act.StartTime,
			&act.EndTime,
			&act.TicketCount,
			&act.InvitationCount,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, act)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}

func (p *PostgresDB) GetActivity(id uuid.UUID) (*models.ActivityWithScannedUser, error) {
	scannedUsers := 0
	a := &models.ActivityWithScannedUser{}
	query := `SELECT id, event_id, name, hall_name, start_time, end_time FROM activities WHERE id = $1 AND deleted_at IS NULL`
	err := p.sql.QueryRow(query, id).Scan(&a.ID, &a.EventID, &a.Name, &a.HallName, &a.StartTime, &a.EndTime)
	if err != nil {
		return nil, err
	}

	a.NumberOfScanedUsers = scannedUsers
	return a, err
}

func (p *PostgresDB) UpdateActivity(a *models.ActivityWithScannedUser) error {
	query := `UPDATE activities SET event_id=$1, name=$2, hall_name=$3, start_time=$4, end_time=$5 WHERE id=$6`
	_, err := p.sql.Exec(query, a.EventID, a.Name, a.HallName, a.StartTime, a.EndTime, a.ID)
	return err
}

func (p *PostgresDB) DeleteActivity(id uuid.UUID) error {
	_, err := p.sql.Exec(`DELETE FROM activities WHERE id=$1`, id)
	return err
}

func (p *PostgresDB) GetActivitiesByEvent(firebaseId string, eventID uuid.UUID) ([]models.ActivityWithScannedUser, error) {
	activities := []models.ActivityWithScannedUser{}

	query := `
SELECT
  a.id,
  a.event_id,
  a.name,
  a.hall_name,
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
		var a models.ActivityWithScannedUser
		if err := rows.Scan(&a.ID, &a.EventID, &a.Name, &a.HallName, &a.StartTime, &a.EndTime, &a.NumberOfScanedUsers); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}

	return activities, nil
}
