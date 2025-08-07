package postgres

import (
	"database/sql"
)


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
