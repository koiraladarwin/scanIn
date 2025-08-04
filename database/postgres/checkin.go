package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	db "github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateCheckInLog(c *models.CheckInLog) error {
	var id string
	query := `INSERT INTO check_in_logs (user_id, activity_id, scanned_at, status, scanned_by)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return p.sql.QueryRow(query, c.UserID, c.ActivityID, c.ScannedAt, c.Status, c.ScannedBy).Scan(&id)
}

func (p *PostgresDB) GetCheckInLog(id uuid.UUID) (*models.CheckInLog, error) {
	c := &models.CheckInLog{}
	query := `SELECT id, user_id, activity_id, scanned_at, status, scanned_by FROM check_in_logs WHERE id=$1`
	err := p.sql.QueryRow(query, id).Scan(&c.ID, &c.UserID, &c.ActivityID, &c.ScannedAt, &c.Status, &c.ScannedBy)
	return c, err
}

func (p *PostgresDB) GetAllCheckInLog() ([]models.CheckInLog, error) {
	logs := []models.CheckInLog{}
	query := `SELECT id, user_id, activity_id, scanned_at, status, scanned_by FROM check_in_logs`

	rows, err := p.sql.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var log models.CheckInLog
		err := rows.Scan(&log.ID, &log.UserID, &log.ActivityID, &log.ScannedAt, &log.Status, &log.ScannedBy)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (p *PostgresDB) UpdateCheckInLog(c *models.CheckInLog) error {
	query := `UPDATE check_in_logs SET user_id=$1, activity_id=$2, scanned_at=$3, status=$4, scanned_by=$5 WHERE id=$6`
	_, err := p.sql.Exec(query, c.UserID, c.ActivityID, c.ScannedAt, c.Status, c.ScannedBy, c.ID)
	return err
}

func (p *PostgresDB) DeleteCheckInLog(id uuid.UUID) error {
	_, err := p.sql.Exec(`DELETE FROM check_in_logs WHERE id=$1`, id)
	return err
}

func (p *PostgresDB) CheckInExists(attendeeID uuid.UUID, activityID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	query := `SELECT id FROM check_in_logs WHERE user_id = $1 AND activity_id = $2`
	err := p.sql.QueryRow(query, attendeeID, activityID).Scan(&id)

	if err == sql.ErrNoRows {
		return uuid.Nil, db.ErrNotFound
	}

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (p *PostgresDB) GetAllCheckInOfEvents(eventID uuid.UUID) ([]models.CheckInLog, error) {
	var checkIns []models.CheckInLog

	queryActivities := `SELECT id FROM activities WHERE event_id = $1`
	rows, err := p.sql.Query(queryActivities, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var activityID uuid.UUID
		if err := rows.Scan(&activityID); err != nil {
			return nil, err
		}

		queryCheckIn := `SELECT id, user_id, activity_id, scanned_at, status, scanned_by FROM check_in_logs WHERE activity_id = $1`
		activityRows, err := p.sql.Query(queryCheckIn, activityID)
		if err != nil {
			return nil, err
		}

		for activityRows.Next() {
			var checkIn models.CheckInLog
			if err := activityRows.Scan(&checkIn.ID, &checkIn.UserID, &checkIn.ActivityID, &checkIn.ScannedAt, &checkIn.Status, &checkIn.ScannedBy); err != nil {
				activityRows.Close()
				return nil, err
			}
			checkIns = append(checkIns, checkIn)
		}
		activityRows.Close()
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return checkIns, nil
}


func (p *PostgresDB) GetAllCheckInOfActivity(activityID uuid.UUID) ([]models.CheckInRespose, error) {
	var checkIns []models.CheckInRespose

	query := `
		SELECT 
			c.id,
			u.full_name,
      u.auto_id,
      u.role,
			c.user_id,
			a.name as activity_name,
			c.activity_id,
			c.scanned_at,
			c.status,
			c.scanned_by
		FROM check_in_logs c
		JOIN users u ON u.id = c.user_id
		JOIN activities a ON a.id = c.activity_id
		WHERE c.activity_id = $1
	`

	rows, err := p.sql.Query(query, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var checkIn models.CheckInRespose
		if err := rows.Scan(
			&checkIn.ID,
			&checkIn.FullName,
      &checkIn.AutoId,
      &checkIn.Role,
			&checkIn.UserID,
			&checkIn.ActivityName,
			&checkIn.ActivityID,
			&checkIn.ScannedAt,
			&checkIn.Status,
			&checkIn.ScannedBy,
		); err != nil {
			return nil, err
		}
		checkIns = append(checkIns, checkIn)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return checkIns, nil
}

