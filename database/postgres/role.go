package postgres

import (
	"database/sql"

	"github.com/koiraladarwin/scanin/models"
)

func (postgres *PostgresDB) AddEventRole(role models.RoleRequest) error {
	query := `INSERT INTO eventRoles (fireBaseId, event_id, canSeeScanned, canCreateAttendee, canSeeAttendee) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := postgres.sql.Exec(query, role.FireBaseId, role.EventId, role.CanSeeScanned, role.CanAddAttendee, role.CanSeeAttendee)
	return err
}
func (postgres *PostgresDB) ModifyEventRole(role models.EditRoleRequest) error {
	query := `UPDATE eventRoles SET canSeeScanned = $1, canCreateAttendee = $2, canSeeAttendee = $3 
        WHERE fireBaseId = $4 AND event_id = $5`
	_, err := postgres.sql.Exec(query, role.CanSeeScanned, role.CanAddAttendee, role.CanSeeAttendee, role.FireBaseId, role.EventId)
	return err
}

func (postgres *PostgresDB) AddStaffToEvent(fbId, eventId string) error {
	query := `INSERT INTO eventRoles (fireBaseId, event_id) VALUES ($1, $2)`
	_, err := postgres.sql.Exec(query, fbId, eventId)
	return err
}

func (postgres *PostgresDB) GetStaffByEvent(eventId string) ([]models.Staff, error) {
	var fireBaseIds []models.Staff

	query := `SELECT fireBaseId,canSeeScanned, canCreateAttendee, canSeeAttendee FROM eventRoles WHERE event_id = $1`
	rows, err := postgres.sql.Query(query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fireBaseId models.Staff
		err := rows.Scan(&fireBaseId.FireBaseId, &fireBaseId.CanSeeScanned, &fireBaseId.CanCreateAttendee, &fireBaseId.CanSeeAttendee)
		if err != nil {
			continue
		}
		fireBaseIds = append(fireBaseIds, fireBaseId)
	}

	return fireBaseIds, nil
}

func (postgres *PostgresDB) AddAdminToEvent(fbId, eventId string) error {
	query := `INSERT INTO eventRoles (fireBaseId, event_id ,isCreator) VALUES ($1, $2, true)`
	_, err := postgres.sql.Exec(query, fbId, eventId)
	return err
}

func (postgres *PostgresDB) IsCreator(fbId, eventId string) (bool, error) {
	var isCreator bool
	query := `SELECT isCreator FROM eventRoles WHERE fireBaseId = $1 AND event_id = $2`
	err := postgres.sql.QueryRow(query, fbId, eventId).Scan(&isCreator)
	return isCreator, err
}

func (postgres *PostgresDB) CanSeeScanned(fbId, eventId string) (bool, error) {
	var isCreator, canSee bool
	query := `SELECT isCreator, canSeeScanned FROM eventRoles WHERE fireBaseId = $1 AND event_id = $2`
	err := postgres.sql.QueryRow(query, fbId, eventId).Scan(&isCreator, &canSee)
	if err != nil {
		return false, err
	}
	if isCreator {
		return true, nil
	}

	return canSee, nil
}

func (postgres *PostgresDB) CanCreateActivity(fbId, eventId string) (bool, error) {
	var isCreator, canCreate bool
	query := `SELECT isCreator, canCreateActivity FROM eventRoles WHERE fireBaseId = $1 AND event_id = $2`
	err := postgres.sql.QueryRow(query, fbId, eventId).Scan(&isCreator, &canCreate)
	if err != nil {
		return false, err
	}
	if isCreator {
		return true, nil
	}
	return canCreate, nil
}

func (postgres *PostgresDB) CanCreateAttendee(fbId, eventId string) (bool, error) {
	var isCreator, canCreate bool
	query := `SELECT isCreator, canCreateAttendee FROM eventRoles WHERE fireBaseId = $1 AND event_id = $2`
	err := postgres.sql.QueryRow(query, fbId, eventId).Scan(&isCreator, &canCreate)
	if err != nil {
		return false, err
	}
	if isCreator {
		return true, nil
	}
	return canCreate, nil
}

func (postgres *PostgresDB) CanSeeAttendee(fbId, eventId string) (bool, error) {
	var isCreator, canSee bool
	query := `SELECT isCreator, canSeeAttendee FROM eventRoles WHERE fireBaseId = $1 AND event_id = $2`
	err := postgres.sql.QueryRow(query, fbId, eventId).Scan(&isCreator, &canSee)
	if err != nil {
		return false, err
	}
	if isCreator {
		return true, nil
	}
	return canSee, nil
}

func (postgres *PostgresDB) CanSeeEventInfo(fbId, eventId string) (bool, error) {
	query := `SELECT 1 FROM eventRoles WHERE fireBaseId = $1 AND event_id = $2 LIMIT 1`
	var exists int
	err := postgres.sql.QueryRow(query, fbId, eventId).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
