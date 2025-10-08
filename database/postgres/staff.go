package postgres

import (
	db "github.com/koiraladarwin/scanin/database"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateStaffCategory(staffCategoryRequest *models.StaffCategory) (*models.StaffCategory, error) {
	var StaffCategory models.StaffCategory
	query := `
  INSERT INTO staff_category (firebase_id, tag, description)
  VALUES ($1, $2, $3)
  RETURNING id
  `
	err := p.sql.QueryRow(
		query,
		staffCategoryRequest.FirebaseID,
		staffCategoryRequest.Tag,
		staffCategoryRequest.Description,
	).Scan(&StaffCategory.ID)

	if isUniqueViolationError(err) {
		return nil, db.ErrAlreadyExists
	}
	StaffCategory.FirebaseID = staffCategoryRequest.FirebaseID
	StaffCategory.Tag = staffCategoryRequest.Tag
	StaffCategory.Description = staffCategoryRequest.Description

	return &StaffCategory, err
}

func (p *PostgresDB) GetStaffCategories(firebaseId string)([]models.StaffCategory,error){
  query := `SELECT id, firebase_id, tag, description FROM staff_category WHERE firebase_id=$1 AND deleted_at IS NULL`
  rows, err := p.sql.Query(query, firebaseId)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  var userCategories []models.StaffCategory
  for rows.Next() {
    var uc models.StaffCategory
    if err := rows.Scan(&uc.ID, &uc.FirebaseID, &uc.Tag, &uc.Description); err != nil {
      return nil, err
    }
    userCategories = append(userCategories, uc)
  }
  return userCategories, nil
}

func(p *PostgresDB) CreateStaff(staffRequest *models.Staff)(*models.Staff,error){
  var staff models.Staff
  query := `
  INSERT INTO staff (firebase_id, staff_gmail, name, image_url, phone, staff_category_id)
  VALUES ($1, $2, $3, $4, $5, $6)
  RETURNING id
  `
  err := p.sql.QueryRow(
    query,
    staffRequest.FirebaseID,
    staffRequest.StaffGmail,
    staffRequest.Name,
    staffRequest.ImageURL,
    staffRequest.Phone,
    staffRequest.StaffCategoryID,
  ).Scan(&staff.ID)
  if isUniqueViolationError(err) {
    return nil, db.ErrAlreadyExists
  }
  staff.FirebaseID = staffRequest.FirebaseID
  staff.StaffGmail = staffRequest.StaffGmail
  staff.Name = staffRequest.Name
  staff.ImageURL = staffRequest.ImageURL
  staff.Phone = staffRequest.Phone
  staff.StaffCategoryID = staffRequest.StaffCategoryID

  return &staff, err
}

func(p *PostgresDB) GetStaffs(firebaseId string)([]models.Staff,error){
  query := `SELECT id, firebase_id, staff_gmail, name, image_url, phone, staff_category_id FROM staff WHERE firebase_id=$1 AND deleted_at IS NULL`
  rows, err := p.sql.Query(query, firebaseId)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  var staffs []models.Staff
  for rows.Next() {
    var s models.Staff
    if err := rows.Scan(&s.ID, &s.FirebaseID, &s.StaffGmail, &s.Name, &s.ImageURL, &s.Phone, &s.StaffCategoryID); err != nil {
      return nil, err
    }
    staffs = append(staffs, s)
  }
  return staffs, nil 
}
