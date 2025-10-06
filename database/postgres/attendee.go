package postgres

import "github.com/koiraladarwin/scanin/models"

func (p *PostgresDB) CreateAttendeeCategory(a models.AttendeeCategory) (models.AttendeeCategory, error) {
	query := `INSERT INTO attendee_category ( firebase_id, tag, description) VALUES ($1, $2, $3) RETURNING id, firebase_id, tag, description;`
	var attendeeCategory models.AttendeeCategory

	err := p.sql.QueryRow(query, a.FirebaseID, a.Tag, a.Description).Scan(&attendeeCategory.ID, &attendeeCategory.FirebaseID, &attendeeCategory.Tag, &attendeeCategory.Description)
	if err != nil {
		return models.AttendeeCategory{}, err
	}

	return attendeeCategory, nil

}

func (p *PostgresDB) GetAttendeeCategories(firebaseId string) ([]models.AttendeeCategory, error) {
  query := `SELECT id, firebase_id, tag, description FROM attendee_category WHERE firebase_id = $1 AND deleted_at IS NULL;`
  rows, err := p.sql.Query(query, firebaseId)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  var categories []models.AttendeeCategory
  for rows.Next() {
    category := models.AttendeeCategory{} 
    err := rows.Scan(&category.ID, &category.FirebaseID, &category.Tag, &category.Description)
    if err != nil {
      return nil, err
    }
    categories = append(categories, category)
  }
  return categories, nil
}

func (p *PostgresDB) CreateAttendee(a models.AttendeeRequest) (models.Attendee, error) {

	query := `INSERT INTO attendee ( user_id, ticket_id,attendee_category_id) VALUES ($1, $2, $3) RETURNING id,  user_id, ticket_id, attendee_category_id;`
	var attendee models.Attendee

	err := p.sql.QueryRow(query, a.UserID, a.TicketID,a.AttendeeCategoryID).Scan(&attendee.ID, &attendee.UserID, &attendee.TicketID,&attendee.AttendeeCategoryID)
	if err != nil {
		return models.Attendee{}, err
	}
	return attendee, nil
}
