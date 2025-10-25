package postgres

import (
	"log"

	"github.com/google/uuid"
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

func (p *PostgresDB) GetStaffCategories(firebaseId string) ([]models.StaffCategory, error) {
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

func (p *PostgresDB) CreateStaff(staffRequest *models.Staff) (*models.Staff, error) {
	staff_count_query := `SELECT COUNT(*) FROM staff WHERE firebase_id=$1 AND deleted_at IS NULL`
	err := p.sql.QueryRow(staff_count_query, staffRequest.FirebaseID).Scan(&staffRequest.AutoID)
	if err != nil {
		return nil, err
	}
	staffRequest.AutoID += 1
	var staff models.Staff
	query := `
  INSERT INTO staff (firebase_id, staff_gmail, name, image_url, phone, staff_category_id, company, position, auto_id)
  VALUES ($1, $2, $3, $4, $5, $6 ,$7, $8, $9)
  RETURNING id
  `
	err = p.sql.QueryRow(
		query,
		staffRequest.FirebaseID,
		staffRequest.StaffGmail,
		staffRequest.Name,
		staffRequest.ImageURL,
		staffRequest.Phone,
		staffRequest.StaffCategoryID,
		staffRequest.Company,
		staffRequest.Position,
		staffRequest.AutoID,
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
	staff.Company = staffRequest.Company
	staff.Position = staffRequest.Position

	return &staff, err
}

func (p *PostgresDB) GetStaffs(firebaseId string) ([]models.Staff, error) {
	query := `SELECT id, firebase_id, staff_gmail, name, image_url, phone, staff_category_id ,company ,position, auto_id FROM staff WHERE firebase_id=$1 AND deleted_at IS NULL`
	rows, err := p.sql.Query(query, firebaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var staffs []models.Staff
	for rows.Next() {
		var s models.Staff
		if err := rows.Scan(&s.ID, &s.FirebaseID, &s.StaffGmail, &s.Name, &s.ImageURL, &s.Phone, &s.StaffCategoryID, &s.Company, &s.Position, &s.AutoID); err != nil {
			return nil, err
		}
		staffs = append(staffs, s)
	}
	return staffs, nil
}

func (p *PostgresDB) CreateStaffEventEnroll(staffEnrollRequest *models.StaffEnroll) (*models.StaffEnroll, error) {
	var staffEnroll models.StaffEnroll
	query := `
  INSERT INTO staff_enroll (firebase_id, staff_id, event_id, active)
  VALUES ($1, $2, $3, $4)
  RETURNING id
  `
	err := p.sql.QueryRow(
		query,
		staffEnrollRequest.FirebaseID,
		staffEnrollRequest.StaffID,
		staffEnrollRequest.EventID,
		staffEnrollRequest.Active,
	).Scan(&staffEnroll.ID)
	if isUniqueViolationError(err) {
		return nil, db.ErrAlreadyExists
	}
	staffEnroll.FirebaseID = staffEnrollRequest.FirebaseID
	staffEnroll.StaffID = staffEnrollRequest.StaffID
	staffEnroll.EventID = staffEnrollRequest.EventID
	staffEnroll.Active = staffEnrollRequest.Active

	return &staffEnroll, err
}

func (p *PostgresDB) GetStaffEventEnroll(firebaseId string, staffId uuid.UUID, eventId uuid.UUID) (models.StaffEnroll, error) {
	query := `SELECT id, firebase_id, staff_id, event_id, active FROM staff_enroll WHERE firebase_id=$1 AND staff_id=$2 AND event_id=$3 AND deleted_at IS NULL`
	var staffEnroll models.StaffEnroll
	err := p.sql.QueryRow(query, firebaseId, staffId, eventId).Scan(&staffEnroll.ID, &staffEnroll.FirebaseID, &staffEnroll.StaffID, &staffEnroll.EventID, &staffEnroll.Active)
  log.Println("staffId:", staffId)
  log.Println("eventId:", eventId)
	if err != nil {
    log.Println("Error retrieving staff enrollment:", err)
		return staffEnroll, err
	}
	return staffEnroll, nil
}

func (p *PostgresDB) CreateStaffActivityAssign(staffActivityRequest *models.StaffActivities) (*models.StaffActivities, error) {
	var staffActivity models.StaffActivities
	query := `
  INSERT INTO staff_activites (firebase_id, staff_enroll_id, activity_id)
  VALUES ($1, $2, $3)
  RETURNING id
  `
	err := p.sql.QueryRow(
		query,
		staffActivityRequest.FirebaseID,
		staffActivityRequest.StaffEnrollId,
		staffActivityRequest.ActivityID,
	).Scan(&staffActivity.ID)
	if isUniqueViolationError(err) {
		return nil, db.ErrAlreadyExists
	}
	staffActivity.FirebaseID = staffActivityRequest.FirebaseID
	staffActivity.StaffEnrollId = staffActivityRequest.StaffEnrollId
	staffActivity.ActivityID = staffActivityRequest.ActivityID

	return &staffActivity, err
}

func (p *PostgresDB) GetEnrolledStaff(firebaseID string) ([]models.EnrolledStaff, error) {
	query := `
	
SELECT 
	u.auto_id AS auto_id,
	uc.tag AS attendee_category_name,
	u.name AS attendee_name,
	u.image_url AS attendee_image,
	e.name AS event_name,
	a.name AS session_name
FROM staff_activites aa
JOIN staff_enroll at ON aa.staff_enroll_id = at.id
JOIN staff u ON u.id = at.staff_id
LEFT JOIN staff_category uc ON uc.id = u.staff_category_id
LEFT JOIN activities a ON a.id = aa.activity_id
LEFT JOIN events e ON e.id = a.event_id;
	`

	rows, err := p.sql.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendees []models.EnrolledStaff

	for rows.Next() {
		var attendee models.EnrolledStaff
		err := rows.Scan(
			&attendee.AutoID,
			&attendee.StafCategoryName,
			&attendee.StaffName,
			&attendee.StaffImage,
			&attendee.EventName,
			&attendee.ActivityName,
		)

		if err != nil {
			return nil, err
		}
		attendees = append(attendees, attendee)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return attendees, nil
}
