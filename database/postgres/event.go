package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
	"github.com/koiraladarwin/scanin/utils"
	"github.com/lib/pq"
)

func (p *PostgresDB) GetEventsWithDetails(firebaseId string) ([]models.EventWithDetails, error) {
	query := `
SELECT
    e.id,
    e.event_category_id,
    ec.tag AS event_category_name,
    e.name,
    e.event_organizer,
    e.description,
    e.start_time,
    e.end_time,
    e.location,
    COALESCE(array_agg(DISTINCT a.name) FILTER (WHERE a.deleted_at IS NULL AND a.name IS NOT NULL), '{}') AS session_names,
    COALESCE(COUNT(DISTINCT t.id) FILTER (WHERE t.price != 0), 0) AS ticket_count,
    COALESCE(COUNT(DISTINCT t.id) FILTER (WHERE t.price = 0), 0) AS invitation_count
FROM events e
LEFT JOIN event_category ec ON e.event_category_id = ec.id
LEFT JOIN activities a ON a.event_id = e.id
LEFT JOIN ticket t ON t.event_id = e.id
WHERE e.firebase_id = $1
  AND e.deleted_at IS NULL
GROUP BY e.id, ec.tag;
`

	rows, err := p.sql.Query(query, firebaseId)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []models.EventWithDetails

	for rows.Next() {
		var ev models.EventWithDetails
		if err := rows.Scan(
			&ev.ID,
			&ev.EventCategoryID,
			&ev.EventCategoryName,
			&ev.Name,
			&ev.EventOrganizer,
			&ev.Description,
			&ev.StartTime,
			&ev.EndTime,
			&ev.Location,
			pq.Array(&ev.SessionNames),
			&ev.TicketCount,
			&ev.InvitationCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		events = append(events, ev)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return events, nil
}


func (p *PostgresDB) CreateEventCategory(e *models.EventCategoryRequest) (models.EventCategory, error) {
	var eventCategoryID uuid.UUID

	query := `
    INSERT INTO event_category (tag, description,firebase_id)
    VALUES ($1, $2, $3)
    RETURNING id;
  `

	err := p.sql.QueryRow(query,
		e.Tag,
		e.Description,
		e.FirebaseID,
	).Scan(&eventCategoryID)

	if err != nil {
		return models.EventCategory{}, fmt.Errorf("failed to insert event category: %w", err)
	}

	createdEventCategory := models.EventCategory{
		ID:          eventCategoryID,
		Tag:         e.Tag,
		Description: e.Description,
	}

	return createdEventCategory, nil
}

func (p *PostgresDB) GetEventCategories(firebaseID string) ([]models.EventCategory, error) {
	query := `
		SELECT id, tag, description 
		FROM event_category 
		WHERE firebase_id = $1 AND deleted_at IS NULL
	`
	rows, err := p.sql.Query(query, firebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query event categories: %w", err)
	}
	defer rows.Close()

	var categories []models.EventCategory

	for rows.Next() {
		var ec models.EventCategory
		if err := rows.Scan(&ec.ID, &ec.Tag, &ec.Description); err != nil {
			return nil, fmt.Errorf("failed to scan event category: %w", err)
		}
		categories = append(categories, ec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over event categories: %w", err)
	}

	return categories, nil
}

func (p *PostgresDB) CreateEvent(e *models.EventCreateRequest) (models.Event, error) {
	var eventID uuid.UUID
	staffCode := utils.RandomString(6)
	adminCode := utils.RandomString(7)

	query := `
  INSERT INTO events (name, event_category_id, description, start_time, end_time, location, staff_code, admin_code,firebase_id,event_organizer)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id;
	`

	err := p.sql.QueryRow(query,
		e.Name,
		e.EventCategoryID,
		e.Description,
		e.StartTime,
		e.EndTime,
		e.Location,
		staffCode,
		adminCode,
		e.FirebaseID,
		e.EventOrganizer,
	).Scan(&eventID)

	if err != nil {
		return models.Event{}, fmt.Errorf("failed to insert event: %w", err)
	}

	// Build and return the created event
	createdEvent := models.Event{
		ID:          eventID,
		Name:        e.Name,
		Description: e.Description,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Location:    e.Location,
		StaffCode:   &staffCode,
		AdminCode:   &adminCode,
	}

	return createdEvent, nil
}

func (p *PostgresDB) GetEventsWithSessions(firebaseId string) ([]models.EventwithSessions, error) {
	query := `
	SELECT
		e.id AS event_id,
		e.event_category_id,
		ec.tag AS event_category_name,
		e.name AS event_name,
		e.event_organizer,
		e.description,
		e.start_time,
		e.end_time,
		e.location,
		a.id AS activity_id,
		a.name AS activity_name,
		a.hall_name,
		a.start_time AS activity_start,
		a.end_time AS activity_end,
		a.firebase_id AS activity_firebase_id
	FROM events e
	LEFT JOIN event_category ec ON e.event_category_id = ec.id
	LEFT JOIN activities a ON a.event_id = e.id AND a.deleted_at IS NULL
	WHERE e.deleted_at IS NULL AND e.firebase_id = $1
	ORDER BY e.start_time DESC;
	`

	rows, err := p.sql.Query(query, firebaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventMap := make(map[string]*models.EventwithSessions)

	for rows.Next() {
		var (
			eventID             uuid.UUID
			eventCategoryID     uuid.UUID
			eventCategoryName   string
			eventName           string
			eventOrganizer      string
			description         string
			startTime           time.Time
			endTime             time.Time
			location            string
			activityID          sql.NullString
			activityName        sql.NullString
			hallName            sql.NullString
			activityStart       sql.NullTime
			activityEnd         sql.NullTime
			activityFirebaseID  sql.NullString
		)

		err := rows.Scan(
			&eventID,
			&eventCategoryID,
			&eventCategoryName,
			&eventName,
			&eventOrganizer,
			&description,
			&startTime,
			&endTime,
			&location,
			&activityID,
			&activityName,
			&hallName,
			&activityStart,
			&activityEnd,
			&activityFirebaseID,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := eventMap[eventID.String()]; !exists {
			eventMap[eventID.String()] = &models.EventwithSessions{
				ID:                eventID,
				EventCategoryID:   eventCategoryID,
				EventCategoryName: eventCategoryName,
				Name:              eventName,
				EventOrganizer:    eventOrganizer,
				Description:       description,
				StartTime:         startTime,
				EndTime:           endTime,
				Location:          location,
				Activity:          []models.Activity{},
			}
		}

		if activityID.Valid {
			eventMap[eventID.String()].Activity = append(eventMap[eventID.String()].Activity, models.Activity{
				ID:         uuid.MustParse(activityID.String),
				EventID:    eventID,
				FirebaseID: activityFirebaseID.String,
				Name:       activityName.String,
				HallName:   hallName.String,
				StartTime:  activityStart.Time,
				EndTime:    activityEnd.Time,
			})
		}
	}

	var results []models.EventwithSessions
	for _, evt := range eventMap {
		results = append(results, *evt)
	}

	return results, nil
}

func (p *PostgresDB) GetEventsByFirebaseUser(firebaseId string) ([]models.Event, error) {
	query := `
SELECT 
    id,
    name,
    event_organizer,
    event_category_id,
    description,
    start_time,
    end_time,
    location,
    staff_code,
    admin_code
FROM events 
  where firebase_id = $1 AND deleted_at IS NULL;
`
	rows, err := p.sql.Query(query, firebaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var e models.Event
		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.EventOrganizer,
			&e.EventCategoryID,
			&e.Description,
			&e.StartTime,
			&e.EndTime,
			&e.Location,
			&e.StaffCode,
			&e.AdminCode,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, e)
		log.Printf("Event fetched: %+v\n", e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}


func (p *PostgresDB) GetEventsWithSessionsAndTickets(firebaseId string) ([]models.EventwithSessionsAndTickets, error) {
	query := `
	SELECT
		e.id AS event_id,
		e.event_category_id,
		ec.tag AS event_category_name,
		e.name AS event_name,
		e.event_organizer,
		e.description,
		e.start_time,
		e.end_time,
		e.location,

		COALESCE(
			ARRAY_AGG(
				DISTINCT a.id || '|' || a.name || '|' || a.hall_name || '|' ||
				COALESCE(a.start_time::text, '') || '|' ||
				COALESCE(a.end_time::text, '') || '|' ||
				COALESCE(a.firebase_id, '')
			) FILTER (WHERE a.id IS NOT NULL),
			'{}'
		) AS activities,

		COALESCE(
			ARRAY_AGG(
				DISTINCT t.id || '|' || t.name || '|' ||
				COALESCE(t.price::text, '0') || '|' ||
				COALESCE(t.paid::text, 'false') || '|' ||
				COALESCE(t.start_time::text, '') || '|' ||
				COALESCE(t.end_time::text, '') || '|' ||
				COALESCE(t.firebase_id, '')
			) FILTER (WHERE t.id IS NOT NULL),
			'{}'
		) AS tickets

	FROM events e
	LEFT JOIN event_category ec ON e.event_category_id = ec.id
	LEFT JOIN activities a ON a.event_id = e.id AND a.deleted_at IS NULL
	LEFT JOIN ticket t ON t.event_id = e.id AND t.deleted_at IS NULL
	WHERE e.deleted_at IS NULL AND e.firebase_id = $1
	GROUP BY
		e.id, e.event_category_id, ec.tag, e.name,
		e.event_organizer, e.description, e.start_time,
		e.end_time, e.location
	ORDER BY e.start_time DESC;
	`

	rows, err := p.sql.Query(query, firebaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.EventwithSessionsAndTickets

	for rows.Next() {
		var (
			event models.EventwithSessionsAndTickets
			activityArr pq.StringArray
			ticketArr   pq.StringArray
		)

		err := rows.Scan(
			&event.ID,
			&event.EventCategoryID,
			&event.EventCategoryName,
			&event.Name,
			&event.EventOrganizer,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.Location,
			&activityArr,
			&ticketArr,
		)
		if err != nil {
			return nil, err
		}

		// Parse activities
		for _, s := range activityArr {
			parts := strings.Split(s, "|")
			if len(parts) >= 6 {
				start, _ := time.Parse(time.RFC3339, parts[3])
				end, _ := time.Parse(time.RFC3339, parts[4])
				event.Activity = append(event.Activity, models.Activity{
					ID:         uuid.MustParse(parts[0]),
					EventID:    event.ID,
					Name:       parts[1],
					HallName:   parts[2],
					StartTime:  start,
					EndTime:    end,
					FirebaseID: parts[5],
				})
			}
		}

		// Parse tickets
		for _, s := range ticketArr {
			parts := strings.Split(s, "|")
			if len(parts) >= 7 {
				start, _ := time.Parse(time.RFC3339, parts[4])
				end, _ := time.Parse(time.RFC3339, parts[5])
				price, _ := strconv.ParseFloat(parts[2], 64)
				paid := parts[3] == "true"
				event.Ticket = append(event.Ticket, models.Ticket{
					ID:         parts[0],
					Name:       parts[1],
					EventID:    event.ID.String(),
					Price:      price,
					Paid:       paid,
					StartTime:  start.String(),
					EndTime:    end.String(),
					FirebaseID: parts[6],
				})
			}
		}

		results = append(results, event)
	}

	return results, nil
}
func (p *PostgresDB) GetEventByFirebaseUser(firebaseId string, eventId uuid.UUID) (*models.Event, error) {
	e := &models.Event{}

	query := `
	SELECT
		e.id,
		e.name,
    e.event_organizer,
		e.description,
		e.start_time,
		e.end_time,
		e.location,
		e.staff_code,
		e.admin_code,
		COUNT(DISTINCT a.id) AS number_of_participant
	FROM events e
	LEFT JOIN attendee a ON e.id = a.event_id AND a.deleted_at IS NULL
	WHERE e.id = $1 
	  AND e.firebase_id = $2
	  AND e.deleted_at IS NULL
	GROUP BY 
		e.id, e.name, e.description, e.start_time, e.end_time, e.location, e.staff_code, e.admin_code;
	`

	err := p.sql.QueryRow(query, eventId, firebaseId).Scan(
		&e.ID,
		&e.Name,
		&e.EventOrganizer,
		&e.Description,
		&e.StartTime,
		&e.EndTime,
		&e.Location,
		&e.StaffCode,
		&e.AdminCode,
		&e.NumberOfParticipant,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	fmt.Printf("Number of participants: %d\n", e.NumberOfParticipant)
	return e, nil
}

func (p *PostgresDB) GetEventByStaffId(id string) (*models.Event, error) {
	e := &models.Event{}
	query := `SELECT id, name, description, start_time, end_time, location FROM events WHERE staff_code = $1 AND delete_at IS NULL`
	err := p.sql.QueryRow(query, id).Scan(&e.ID, &e.Name, &e.Description, &e.StartTime, &e.EndTime, &e.Location)
	return e, err
}

func (p *PostgresDB) GetEventByAdminId(id string) (*models.Event, error) {
	e := &models.Event{}
	query := `SELECT id, name, description, start_time, end_time, location FROM events WHERE admin_code = $1 And delete_at IS NULL`
	err := p.sql.QueryRow(query, id).Scan(&e.ID, &e.Name, &e.Description, &e.StartTime, &e.EndTime, &e.Location)
	return e, err
}

func (p *PostgresDB) UpdateEvent(e *models.EventModifyRequest) error {
	query := `UPDATE events SET name=$1, description=$2,location=$3 WHERE id=$4`
	_, err := p.sql.Exec(query, e.Name, e.Description, e.Location, e.ID)
	return err
}

func (p *PostgresDB) DeleteEvent(id uuid.UUID) error {
	_, err := p.sql.Exec(`DELETE FROM events WHERE id=$1`, id)
	return err
}

func (p *PostgresDB) EventExists(eventID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1 AND delete_at IS NULL)`
	err := p.sql.QueryRow(query, eventID).Scan(&exists)
	return exists, err
}

func (p *PostgresDB) GetAllEvents() ([]models.Event, error) {
	query := `SELECT id, name, description, start_time, end_time, location FROM events ORDER BY start_time WHERE delete_at IS NULL`

	rows, err := p.sql.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.StartTime, &e.EndTime, &e.Location); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (p *PostgresDB) GetEventIdByActivity(acitvity uuid.UUID) (uuid.UUID, error) {
	event := models.Event{}

	query := `SELECT event_id FROM activities WHERE id = $1 limit 1`
	err := p.sql.QueryRow(query, acitvity).Scan(&event.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get event id by activity: %w", err)
	}
	if event.ID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("no event found for activity id: %s", acitvity)
	}

	return event.ID, nil
}
