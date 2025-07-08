package postgres

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateUser(u *models.User) error {
	query := `INSERT INTO users (full_name, email, phone) VALUES ($1, $2, $3) RETURNING id`
	err := p.sql.QueryRow(query, u.FullName, u.Email, u.Phone).Scan(&u.ID)
	if isUniqueViolationError(err) {
		return db.ErrAlreadyExists
	}
  return err
}

func (p *PostgresDB) GetUser(id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, full_name, email, phone FROM users WHERE id=$1`
	err := p.sql.QueryRow(query, id).Scan(&u.ID, &u.FullName, &u.Email, &u.Phone)
	return u, err
}

func (p *PostgresDB) UpdateUser(u *models.User) error {
	query := `UPDATE users SET full_name=$1, email=$2, phone=$3 WHERE id=$5`
	_, err := p.sql.Exec(query, u.FullName, u.Email, u.Phone, u.ID)
	return err
}

func (p *PostgresDB) DeleteUser(id uuid.UUID) error {
	_, err := p.sql.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

func (p *PostgresDB) GetUsersByEvent(eventID uuid.UUID) ([]models.User, error) {
	var users []models.User

	rows, err := p.sql.Query(`
		SELECT u.id, u.full_name, u.email, u.phone
		FROM attendees a
		JOIN users u ON u.id = a.user_id
		WHERE a.event_id = $1
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Phone); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
