package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	db "github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
)

// CreateCheckInLog inserts a new check-in log and returns the generated UUID
func (p *PostgresDB) CreateCheckInLog(c *models.CheckInLog) error {
	query := `INSERT INTO check_in_logs (attendee_id, activity_id, scanned_at, status, scanned_by)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return p.sql.QueryRow(query, c.AttendeeID, c.ActivityID, c.ScannedAt, c.Status, c.ScannedBy).Scan(&c.ID)
}

// GetCheckInLog fetches a check-in log by UUID
func (p *PostgresDB) GetCheckInLog(id uuid.UUID) (*models.CheckInLog, error) {
	c := &models.CheckInLog{}
	query := `SELECT id, attendee_id, activity_id, scanned_at, status, scanned_by FROM check_in_logs WHERE id=$1`
	err := p.sql.QueryRow(query, id).Scan(&c.ID, &c.AttendeeID, &c.ActivityID, &c.ScannedAt, &c.Status, &c.ScannedBy)
	return c, err
}


func (p *PostgresDB) GetAllCheckInLog() ([]models.CheckInLog, error) {
    logs := []models.CheckInLog{}
    query := `SELECT id, attendee_id, activity_id, scanned_at, status, scanned_by FROM check_in_logs`

    rows, err := p.sql.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var log models.CheckInLog
        err := rows.Scan(&log.ID, &log.AttendeeID, &log.ActivityID, &log.ScannedAt, &log.Status, &log.ScannedBy)
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


// UpdateCheckInLog updates an existing check-in log
func (p *PostgresDB) UpdateCheckInLog(c *models.CheckInLog) error {
	query := `UPDATE check_in_logs SET attendee_id=$1, activity_id=$2, scanned_at=$3, status=$4, scanned_by=$5 WHERE id=$6`
	_, err := p.sql.Exec(query, c.AttendeeID, c.ActivityID, c.ScannedAt, c.Status, c.ScannedBy, c.ID)
	return err
}

// DeleteCheckInLog deletes a check-in log by UUID
func (p *PostgresDB) DeleteCheckInLog(id uuid.UUID) error {
	_, err := p.sql.Exec(`DELETE FROM check_in_logs WHERE id=$1`, id)
	return err
}

func (p *PostgresDB) CheckInExists(attendeeID uuid.UUID, activityID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	query := `SELECT id FROM check_in_logs WHERE attendee_id = $1 AND activity_id = $2`
	err := p.sql.QueryRow(query, attendeeID, activityID).Scan(&id)

	if err == sql.ErrNoRows {
		return uuid.Nil, db.ErrNotFound
	}

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
