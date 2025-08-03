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

		`CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			description TEXT,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			location TEXT
		);`,

		`create table if not exists users (
			id uuid primary key default gen_random_uuid(),
      auto_id int not null,
			full_name text not null,
      image_url text not null,
			company text not null,
			position text not null,
			role text not null,
			event_id UUID NOT NULL REFERENCES events(id),
      unique(role,auto_id,event_id)
		);`,

		`CREATE TABLE IF NOT EXISTS activities (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS staff (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			full_name TEXT NOT NULL,
      password TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			phone TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS check_in_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
			scanned_at TIMESTAMP NOT NULL DEFAULT now(),
			status TEXT NOT NULL,
			scanned_by TEXT NOT NULL ,
      UNIQUE (user_id, activity_id)
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
