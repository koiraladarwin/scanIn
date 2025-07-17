package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/koiraladarwin/scanin/database"
)

type PostgresDB struct {
	sql *sql.DB
}

func ConnectPostgres(connStr string) (db.Database, error) {

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	p := &PostgresDB{sql: db}
	if err := p.createTables(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *PostgresDB) Close() error {
	return p.sql.Close()
}

func (p *PostgresDB) createTables() error {
	stmts := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`, // enables gen_random_uuid() in postgress sometimes it may not be active

		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			full_name TEXT NOT NULL,
      image_url TEXT NOT NULL,
			company TEXT NOT NULL,
			position TEXT NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			description TEXT,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			location TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS activities (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS attendees (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			role TEXT NOT NULL CHECK (role IN ('participant', 'staff', 'member')),
      UNIQUE (user_id, event_id)
		);`, //attendees are users in a event

		`CREATE TABLE IF NOT EXISTS staff (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			full_name TEXT NOT NULL,
      password TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			phone TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS check_in_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			attendee_id UUID NOT NULL REFERENCES attendees(id) ON DELETE CASCADE,
			activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
			scanned_at TIMESTAMP NOT NULL DEFAULT now(),
			status TEXT NOT NULL,
			scanned_by TEXT NOT NULL ,
      UNIQUE (attendee_id, activity_id)
		);`,
	}

	for _, stmt := range stmts {
		if _, err := p.sql.Exec(stmt); err != nil {
			return fmt.Errorf("error running table creation: %w", err)
		}
	}
	return nil
}

func isUniqueViolationError(err error) bool {
	if err == nil {
		return false
	}
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}
	return pgErr.Code == "23505"
}
