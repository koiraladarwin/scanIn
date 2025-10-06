package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateUserCategory(reqUserCat *models.UsersCategoryRequest) (*models.UsersCategory, error) {
	var userCat models.UsersCategory
	query := `
  INSERT INTO users_category (firebase_id, tag, description)
  VALUES ($1, $2, $3)
  RETURNING id
  `
	err := p.sql.QueryRow(
		query,
		reqUserCat.FirebaseID,
		reqUserCat.Tag,
		reqUserCat.Description,
	).Scan(&userCat.ID)

	if isUniqueViolationError(err) {
		return nil, db.ErrAlreadyExists
	}
	userCat.FirebaseID = reqUserCat.FirebaseID
	userCat.Tag = reqUserCat.Tag
	userCat.Description = reqUserCat.Description

	return &userCat, err
}

func (p *PostgresDB) CreateUser(reqUser *models.UserRequest) (*models.User, error) {
	var lastAutoID int
	var user models.User

	err := p.sql.QueryRow(`SELECT COALESCE(MAX(auto_id), 0) FROM users WHERE users_category_id = $1`, reqUser.UsersCategoryID).Scan(&lastAutoID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest auto_id: %w", err)
	}

	autoId := lastAutoID + 1

	query := `
		INSERT INTO users (auto_id, full_name, image_url, position, company, users_category_id,firebase_id)
		VALUES ($1, $2, $3, $4, $5, $6,$7)
		RETURNING id
	`
	err = p.sql.QueryRow(
		query,
		autoId,
		reqUser.FullName,
		reqUser.Image_url,
		reqUser.Position,
		reqUser.Company,
		reqUser.UsersCategoryID,
    reqUser.FirebaseID,
	).Scan(&user.ID)

	if isUniqueViolationError(err) {
		return nil, db.ErrAlreadyExists
	}
	user.FullName = reqUser.FullName
	user.Company = reqUser.Company
	user.Position = reqUser.Position
	user.Image_url = reqUser.Image_url
	user.AutoId = autoId
	user.UsersCategoryID = reqUser.UsersCategoryID

	return &user, err
}

func (p *PostgresDB) GetUsers(firebaseid string) ([]models.User, error) {
	var users []models.User
	query := `SELECT id, full_name, auto_id, image_url, position, company ,users_category_id FROM users WHERE firebase_id=$1 AND deleted_at IS NULL`
	rows, err := p.sql.Query(query, firebaseid)
	if err != nil {
		fmt.Println("Error fetching users:", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.FullName, &u.AutoId, &u.Image_url, &u.Position, &u.Company, &u.UsersCategoryID); err != nil {
			fmt.Println("Error scanning user:", err)
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (p *PostgresDB) UpdateUser(u *models.UserModifyRequest) error {
	query := `UPDATE users SET full_name=$1, image_url=$2, position=$3, company=$4,users_category_id=$5 WHERE id=$6`
	_, err := p.sql.Exec(query, u.FullName, u.Image_url, u.Position, u.Company, u.UsersCategoryID, u.ID)
	return err
}

func (p *PostgresDB) DeleteUser(id uuid.UUID) error {
	query := `UPDATE users SET delete_at=NOW() WHERE id=$1`
	_, err := p.sql.Exec(query, id)
	return err
}

func (p *PostgresDB) GetUser(id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, full_name, auto_id, image_url, position, company ,users_category_id FROM users WHERE id=$1 AND delete_at IS NULL`
	err := p.sql.QueryRow(query, id).Scan(&u.ID, &u.FullName, &u.AutoId, &u.Image_url, &u.Position, &u.Company, &u.UsersCategoryID)
	return u, err
}

func (p *PostgresDB) GetUsersByEvent(eventID uuid.UUID) ([]models.User, error) {
	var users []models.User

	rows, err := p.sql.Query(`
			SELECT id, full_name, auto_id, image_url, position, company ,users_category_id FROM users WHERE event_id = $1 AND delete_at IS NULL
	`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.FullName, &u.AutoId, &u.Image_url, &u.Position, &u.Company, &u.UsersCategoryID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (p *PostgresDB) GetNumberOfUsersByEvent(eventID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE event_id = $1 AND delete_at IS NULL`
	err := p.sql.QueryRow(query, eventID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
