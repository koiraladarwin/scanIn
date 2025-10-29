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
     firebase_id TEXT NOT NULL,
     tag TEXT NOT NULL,
     description TEXT NOT NULL,
     deleted_at TIMESTAMPTZ
    );`,

		// 1. Events table
		`CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      firebase_id TEXT NOT NULL,
      event_category_id UUID REFERENCES event_category(id) ON DELETE SET NULL,
			name TEXT NOT NULL,
      event_organizer TEXT NOT NULL,
			description TEXT NOT NULL,
			start_time TIMESTAMPTZ NOT NULL,
			end_time TIMESTAMPTZ NOT NULL,
			location TEXT NOT NULL,
			admin_code TEXT NOT NULL UNIQUE,
			staff_code TEXT NOT NULL UNIQUE,
			deleted_at TIMESTAMPTZ
		);`,

		// 2. Users Category table
		`CREATE TABLE IF NOT EXISTS users_category(
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
     firebase_id TEXT NOT NULL,
     tag TEXT NOT NULL,
     description TEXT NOT NULL, 
     deleted_at TIMESTAMPTZ
    );`,

		// 2. Users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			auto_id INT NOT NULL,
      firebase_id TEXT NOT NULL,
      phone_number TEXT NOT NULL,
      gmail TEXT NOT NULL,
			full_name TEXT NOT NULL,
			image_url TEXT NOT NULL,
			company TEXT NOT NULL,
			position TEXT NOT NULL,
      users_category_id UUID REFERENCES users_category(id) ON DELETE SET NULL,
			UNIQUE (gmail,firebase_id),
			UNIQUE (users_category_id, auto_id),
			deleted_at TIMESTAMPTZ
		);`,

		// 3. Ticket table category
		`CREATE TABLE IF NOT EXISTS ticket_category(
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      firebase_id TEXT NOT NULL,
      tag TEXT NOT NULL,
      description TEXT NOT NULL,
      type TEXT NOT NULL CHECK (type IN ('inv', 'tkt')),
      deleted_at TIMESTAMPTZ
      );
    `,

		// 3. Ticket table
		`CREATE TABLE IF NOT EXISTS ticket (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      firebase_id TEXT NOT NULL,
      ticket_category_id UUID REFERENCES ticket_category(id) ON DELETE SET NULL,
      name TEXT NOT NULL,
			event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			price NUMERIC NOT NULL,
			start_time TIMESTAMPTZ NOT NULL,
			end_time TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ
		);`,

		// 4. Activities table
		`CREATE TABLE IF NOT EXISTS activities (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			hall_name TEXT NOT NULL,
			start_time TIMESTAMPTZ NOT NULL,
			end_time TIMESTAMPTZ NOT NULL,
			deleted_at TIMESTAMPTZ,
      firebase_id TEXT NOT NULL
		);`,

		// 5. Attendee table
		`CREATE TABLE IF NOT EXISTS attendee (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			ticket_id UUID REFERENCES ticket(id) ON DELETE SET NULL,
      paid BOOLEAN NOT NULL DEFAULT FALSE,
			deleted_at TIMESTAMPTZ,
      UNIQUE (user_id, ticket_id)
		);`,

		// 5. Attendee table
		`CREATE TABLE IF NOT EXISTS attendee_activity (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      firebase_id TEXT NOT NULL,
      attendee_id UUID NOT NULL REFERENCES attendee(id) ON DELETE CASCADE,
      activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE, 
      UNIQUE (attendee_id, activity_id),
			deleted_at TIMESTAMPTZ
		);`,

		// 6. Check-in logs table
		`CREATE TABLE IF NOT EXISTS check_in_logs (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     attendee_id UUID NOT NULL REFERENCES attendee(id) ON DELETE CASCADE,
     activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
     scanned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
     scanned_by TEXT NOT NULL,
     deleted_at TIMESTAMPTZ
      );

      CREATE UNIQUE INDEX IF NOT EXISTS unique_active_checkin
      ON check_in_logs (attendee_id, activity_id)
      WHERE deleted_at IS NULL;`,

		// 7.Staff category table
		`CREATE TABLE IF NOT EXISTS staff_category (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      firebase_id TEXT NOT NULL,
      tag TEXT NOT NULL,
      description TEXT NOT NULL,
      deleted_at TIMESTAMPTZ
    );
    `,

		// 7.Staff table
		`CREATE TABLE IF NOT EXISTS staff (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      firebase_id TEXT NOT NULL,
      auto_id INT NOT NULL,
      staff_gmail TEXT NOT NULL UNIQUE,
      name TEXT NOT NULL,
      image_url TEXT NOT NULL,
      phone TEXT NOT NULl,
      company TEXT NOT NULL,
      position TEXT NOT NULL, 
      staff_category_id UUID REFERENCES staff_category(id) ON DELETE SET NULL,
      deleted_at TIMESTAMPTZ
    );
    `,
		// 7.Staff enroll table
		`CREATE TABLE IF NOT EXISTS staff_enroll (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      firebase_id TEXT NOT NULL,
      staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
      event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
      active BOOLEAN NOT NULL DEFAULT FALSE,
      deleted_at TIMESTAMPTZ,
      UNIQUE (staff_id, event_id)
    );`,
		// 7. Staff activities table
		`CREATE TABLE IF NOT EXISTS staff_activites (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      firebase_id TEXT NOT NULL,
      staff_enroll_id UUID NOT NULL REFERENCES staff_enroll(id) ON DELETE CASCADE,
      activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
      deleted_at TIMESTAMPTZ,
      UNIQUE (staff_enroll_id, activity_id)
    );
    `,
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
