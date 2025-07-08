package postgres

import (
	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateUser(u *models.User) error {
	query := `INSERT INTO users (full_name, email, phone, role) VALUES ($1, $2, $3, $4) RETURNING id`
	return p.sql.QueryRow(query, u.FullName, u.Email, u.Phone, u.Role).Scan(&u.ID)
}

func (p *PostgresDB) GetUser(id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, full_name, email, phone, role FROM users WHERE id=$1`
	err := p.sql.QueryRow(query, id).Scan(&u.ID, &u.FullName, &u.Email, &u.Phone, &u.Role)
	return u, err
}

func (p *PostgresDB) UpdateUser(u *models.User) error {
	query := `UPDATE users SET full_name=$1, email=$2, phone=$3, role=$4 WHERE id=$5`
	_, err := p.sql.Exec(query, u.FullName, u.Email, u.Phone, u.Role, u.ID)
	return err
}

func (p *PostgresDB) DeleteUser(id uuid.UUID) error {
	_, err := p.sql.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

