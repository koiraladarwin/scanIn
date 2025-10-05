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
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`,

		// 1. Events category table
		`CREATE TABLE IF NOT EXISTS event_category(
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
     firebase_id TEXT,
     tag TEXT NOT NULL,
     description TEXT,
     deleted_at TIMESTAMPTZ
    );`,

		// 1. Events table
		`CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      firebase_id TEXT,
      event_category_id UUID REFERENCES event_category(id) ON DELETE SET NULL,
			name TEXT NOT NULL,
			description TEXT,
			start_time TIMESTAMPTZ NOT NULL,
			end_time TIMESTAMPTZ NOT NULL,
			location TEXT,
			admin_code TEXT NOT NULL UNIQUE,
			staff_code TEXT NOT NULL UNIQUE,
			deleted_at TIMESTAMPTZ
		);`,

		// 2. Users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			auto_id INT NOT NULL,
			full_name TEXT NOT NULL,
			image_url TEXT NOT NULL,
			company TEXT NOT NULL,
			position TEXT NOT NULL,
			role TEXT NOT NULL,
			event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			UNIQUE (role, auto_id, event_id),
			deleted_at TIMESTAMPTZ
		);`,

		// 3. Ticket table (must come before attendee)
		`CREATE TABLE IF NOT EXISTS ticket (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			price NUMERIC NOT NULL,
			deleted_at TIMESTAMPTZ
		);`,

		// 4. Attendee table
		`CREATE TABLE IF NOT EXISTS attendee (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			ticket_id UUID REFERENCES ticket(id) ON DELETE SET NULL,
			deleted_at TIMESTAMPTZ
		);`,

		// 5. Activities table
		`CREATE TABLE IF NOT EXISTS activities (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			start_time TIMESTAMPTZ NOT NULL,
			end_time TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ
		);`,

		// 7. Scan roles table
		`CREATE TABLE IF NOT EXISTS scan_roles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			firebase_id TEXT NOT NULL,
			activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
			access BOOLEAN NOT NULL DEFAULT false,
			UNIQUE (firebase_id, activity_id)
		);`,

		// 8. Check-in logs table
		`CREATE TABLE IF NOT EXISTS check_in_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			attendee_id UUID NOT NULL REFERENCES attendee(id) ON DELETE CASCADE,
			activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
			scanned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			status TEXT NOT NULL,
			scanned_by TEXT NOT NULL,
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
