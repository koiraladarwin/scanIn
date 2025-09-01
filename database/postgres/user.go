package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateUser(reqUser *models.UserRequest) (*models.User, error) {
	var lastAutoID int
	var user models.User

	err := p.sql.QueryRow(`SELECT COALESCE(MAX(auto_id), 0) FROM users WHERE role = $1`, reqUser.Role).Scan(&lastAutoID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest auto_id: %w", err)
	}

	autoId := lastAutoID + 1

	query := `
		INSERT INTO users (auto_id, full_name, image_url, position, company, role,event_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err = p.sql.QueryRow(
		query,
		autoId,
		reqUser.FullName,
		reqUser.Image_url,
		reqUser.Position,
		reqUser.Company,
		reqUser.Role,
		reqUser.EventId,
	).Scan(&user.ID)

	if isUniqueViolationError(err) {
		return nil, db.ErrAlreadyExists
	}
	user.FullName = reqUser.FullName
	user.Company = reqUser.Company
	user.Position = reqUser.Position
	user.Image_url = reqUser.Image_url
	user.AutoId = autoId
	user.Role = reqUser.Role
	user.EventId = reqUser.EventId

	return &user, err
}

func (p *PostgresDB) GetUser(id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, full_name, auto_id, image_url, position, company ,role,event_id FROM users WHERE id=$1 WHERE delete_at IS NULL`
	err := p.sql.QueryRow(query, id).Scan(&u.ID, &u.FullName, &u.AutoId, &u.Image_url, &u.Position, &u.Company, &u.Role, &u.EventId)
	return u, err
}

func (p *PostgresDB) GetUsersByEvent(eventID uuid.UUID) ([]models.User, error) {
	var users []models.User

	rows, err := p.sql.Query(`
			SELECT id, full_name, auto_id, image_url, position, company ,role,event_id FROM users WHERE event_id = $1 WHERE delete_at IS NULL
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.FullName, &u.AutoId, &u.Image_url, &u.Position, &u.Company, &u.Role, &u.EventId); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (p *PostgresDB) GetNumberOfUsersByEvent(eventID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE event_id = $1 WHERE delete_at IS NULL`
	err := p.sql.QueryRow(query, eventID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
