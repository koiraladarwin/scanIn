package postgres

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/koiraladarwin/scanin/models"
)

func (p *PostgresDB) CreateTicketCategory(a models.TicketCategory) (models.TicketCategory, error) {
	log.Print("Creating ticket category:", a.Type)
	query := `INSERT INTO ticket_category (tag, description, type, firebase_id) VALUES ($1, $2, $3, $4) RETURNING id, tag, description, type;`
	var ticketCategory models.TicketCategory
	err := p.sql.QueryRow(query, a.Tag, a.Description, a.Type, a.FirebaseID).Scan(&ticketCategory.ID, &ticketCategory.Tag, &ticketCategory.Description, &ticketCategory.Type)
	if err != nil {
		return models.TicketCategory{}, err
	}
	return ticketCategory, nil
}

func (p *PostgresDB) GetTicketCategory(firebaseId string, id uuid.UUID) (models.TicketCategory, error) {
	query := `SELECT id, tag, description, type FROM ticket_category WHERE firebase_id = $1 AND deleted_at IS NULL AND id = $2;`
	var tc models.TicketCategory
	err := p.sql.QueryRow(query, firebaseId, id).Scan(&tc.ID, &tc.Tag, &tc.Description, &tc.Type)
	if err != nil {
		return models.TicketCategory{}, err
	}
	return tc, nil
}

func (p *PostgresDB) GetTicketCategories(firebaseId string, ticket_type string) ([]models.TicketCategory, error) {
	query := `SELECT id, tag, description, type FROM ticket_category WHERE firebase_id = $1 AND deleted_at IS NULL AND type=$2;`
	rows, err := p.sql.Query(query, firebaseId, ticket_type)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ticketCategories []models.TicketCategory
	for rows.Next() {
		var tc models.TicketCategory
		if err := rows.Scan(&tc.ID, &tc.Tag, &tc.Description, &tc.Type); err != nil {
			return nil, err
		}
		ticketCategories = append(ticketCategories, tc)
	}
	return ticketCategories, nil
}

func (p *PostgresDB) CreateTicket(a models.TicketRequest) (models.Ticket, error) {
	query := `
	INSERT INTO ticket (event_id, ticket_category_id, price, name, firebase_id, paid, start_time, end_time)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, event_id, ticket_category_id, price, name, firebase_id, paid, start_time, end_time;
	`

	var ticket models.Ticket
	err := p.sql.QueryRow(query, a.EventID, a.TicketCategoryID, a.Price, a.Name, a.FirebaseID, a.Paid, a.StartTime, a.EndTime).
		Scan(&ticket.ID, &ticket.EventID, &ticket.TicketCategoryID, &ticket.Price, &ticket.Name, &ticket.FirebaseID, &ticket.Paid, &ticket.StartTime, &ticket.EndTime)

	if err != nil {
		return models.Ticket{}, err
	}
	return ticket, nil
}

func (p *PostgresDB) GetTickets(firebaseID string) ([]models.Ticket, error) {
	query := `SELECT id, event_id, ticket_category_id, price, name, paid ,start_time, end_time FROM ticket WHERE firebase_id = $1 AND deleted_at IS NULL;`
	rows, err := p.sql.Query(query, firebaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.EventID, &t.TicketCategoryID, &t.Price, &t.Name, &t.Paid, &t.StartTime, &t.EndTime); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func (p *PostgresDB) GetTicketsForAllEvents(firebaseId string) ([]models.EventsTicket, error) {
	query := `
	SELECT
		e.id AS event_id,
		e.name AS event_name,
		e.description AS event_description,
		e.start_time AS event_start_time,
		e.end_time AS event_end_time,
		e.event_organizer AS event_organizer,
		e.location AS event_location,

		t.id AS ticket_id,
		t.ticket_category_id,
		tc.tag AS ticket_category_tag,
		t.price,

		a.id AS attendee_id,
		a.user_id AS attendee_user_id,
		a.ticket_id AS attendee_ticket_id,
		a.deleted_at AS attendee_deleted_at

	FROM events e
	LEFT JOIN ticket t ON e.id = t.event_id AND t.deleted_at IS NULL AND t.firebase_id = $1 AND t.price > 0
	LEFT JOIN ticket_category tc ON t.ticket_category_id = tc.id
	LEFT JOIN attendee a ON t.id = a.ticket_id
	WHERE e.deleted_at IS NULL AND e.firebase_id = $1
	ORDER BY e.start_time DESC;
	`

	rows, err := p.sql.Query(query, firebaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventMap := make(map[string]*models.EventsTicket)
	ticketMap := make(map[string]*models.TicketDetails) // key: ticket ID

	for rows.Next() {
		var (
			eventID uuid.UUID
			e       models.EventResponse

			ticketID          sql.NullString
			ticketCategoryID  sql.NullString
			ticketCategoryTag sql.NullString
			price             sql.NullFloat64

			attendeeID        sql.NullString
			attendeeUserID    uuid.NullUUID
			attendeeTicketID  uuid.NullUUID
			attendeeDeletedAt sql.NullString
		)

		err := rows.Scan(
			&eventID,
			&e.Name,
			&e.Description,
			&e.StartTime,
			&e.EndTime,
			&e.EventOrganizer,
			&e.Location,
			&ticketID,
			&ticketCategoryID,
			&ticketCategoryTag,
			&price,
			&attendeeID,
			&attendeeUserID,
			&attendeeTicketID,
			&attendeeDeletedAt,
		)
		if err != nil {
			return nil, err
		}

		e.ID = eventID

		var t *models.TicketDetails
		if ticketID.Valid {
			if existingTicket, ok := ticketMap[ticketID.String]; ok {
				t = existingTicket
			} else {
				newTicket := &models.TicketDetails{
					Id:                  ticketID.String,
					Ticket_category_id:  ticketCategoryID.String,
					Ticket_category_tag: ticketCategoryTag.String,
					Price:               price.Float64,
					Attendees:           []models.Attendee{},
				}
				ticketMap[ticketID.String] = newTicket
				t = newTicket
			}

			if attendeeID.Valid {
				att := models.Attendee{
					ID:       attendeeID.String,
					UserID:   attendeeUserID.UUID,
					TicketID: attendeeTicketID.UUID,
				}
				if attendeeDeletedAt.Valid {
					att.DeletedAt = &attendeeDeletedAt.String
				}
				t.Attendees = append(t.Attendees, att)
			}
		}

		if existingEvent, ok := eventMap[eventID.String()]; ok {
			if t != nil {
				existingEvent.TicketDetails = append(existingEvent.TicketDetails, *t)
			}
		} else {
			newEvent := &models.EventsTicket{
				Event: e,
			}
			if t != nil {
				newEvent.TicketDetails = []models.TicketDetails{*t}
			} else {
				newEvent.TicketDetails = []models.TicketDetails{}
			}
			eventMap[eventID.String()] = newEvent
		}
	}

	var results []models.EventsTicket
	for _, evt := range eventMap {
		results = append(results, *evt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (p *PostgresDB) GetInviteeForAllEvents(firebaseId string) ([]models.EventsInvitee, error) {
	query := `
	SELECT
		e.id AS event_id,
		e.name AS event_name,
		e.description AS event_description,
		e.start_time AS event_start_time,
		e.end_time AS event_end_time,
		e.event_organizer AS event_organizer,
		e.location AS event_location,

		t.id AS ticket_id,
		t.ticket_category_id,
		tc.tag AS ticket_category_tag,
		t.price,

		a.id AS attendee_id,
		a.user_id AS attendee_user_id,
		a.ticket_id AS attendee_ticket_id,
		a.deleted_at AS attendee_deleted_at

	FROM events e
	LEFT JOIN ticket t ON e.id = t.event_id AND t.deleted_at IS NULL AND t.firebase_id = $1 AND t.price = 0
	LEFT JOIN ticket_category tc ON t.ticket_category_id = tc.id
	LEFT JOIN attendee a ON t.id = a.ticket_id
	WHERE e.deleted_at IS NULL AND e.firebase_id = $1
	ORDER BY e.start_time DESC;
	`

	rows, err := p.sql.Query(query, firebaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventMap := make(map[string]*models.EventsInvitee)
	ticketMap := make(map[string]*models.InviteeDetails)

	for rows.Next() {
		var (
			eventID uuid.UUID
			e       models.EventResponse

			ticketID          sql.NullString
			ticketCategoryID  sql.NullString
			ticketCategoryTag sql.NullString
			price             sql.NullFloat64

			attendeeID        sql.NullString
			attendeeUserID    uuid.NullUUID
			attendeeTicketID  uuid.NullUUID
			attendeeDeletedAt sql.NullString
		)

		err := rows.Scan(
			&eventID,
			&e.Name,
			&e.Description,
			&e.StartTime,
			&e.EndTime,
			&e.EventOrganizer,
			&e.Location,
			&ticketID,
			&ticketCategoryID,
			&ticketCategoryTag,
			&price,
			&attendeeID,
			&attendeeUserID,
			&attendeeTicketID,
			&attendeeDeletedAt,
		)
		if err != nil {
			return nil, err
		}

		e.ID = eventID

		var t *models.InviteeDetails
		if ticketID.Valid {
			if existingTicket, ok := ticketMap[ticketID.String]; ok {
				t = existingTicket
			} else {
				newTicket := &models.InviteeDetails{
					Id:                  ticketID.String,
					Ticket_category_id:  ticketCategoryID.String,
					Ticket_category_tag: ticketCategoryTag.String,
					Attendees:           []models.Attendee{},
				}
				ticketMap[ticketID.String] = newTicket
				t = newTicket
			}

			if attendeeID.Valid {
				att := models.Attendee{
					ID:       attendeeID.String,
					UserID:   attendeeUserID.UUID,
					TicketID: attendeeTicketID.UUID,
				}
				if attendeeDeletedAt.Valid {
					att.DeletedAt = &attendeeDeletedAt.String
				}
				t.Attendees = append(t.Attendees, att)
			}
		}

		if existingEvent, ok := eventMap[eventID.String()]; ok {
			if t != nil {
				existingEvent.InviteeDetails = append(existingEvent.InviteeDetails, *t)
			}
		} else {
			newEvent := &models.EventsInvitee{
				Event: e,
			}
			if t != nil {
				newEvent.InviteeDetails = []models.InviteeDetails{*t}
			} else {
				newEvent.InviteeDetails = []models.InviteeDetails{}
			}
			eventMap[eventID.String()] = newEvent
		}
	}

	var results []models.EventsInvitee
	for _, evt := range eventMap {
		results = append(results, *evt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
func (p *PostgresDB) GetTicketById(firebaseId string, id uuid.UUID) (models.Ticket, error) {
	query := `SELECT id, event_id, ticket_category_id, price, name, paid ,start_time, end_time FROM ticket WHERE firebase_id = $1 AND deleted_at IS NULL AND id = $2;`
	var t models.Ticket
	err := p.sql.QueryRow(query, firebaseId, id).Scan(&t.ID, &t.EventID, &t.TicketCategoryID, &t.Price, &t.Name, &t.Paid, &t.StartTime, &t.EndTime)
	if err != nil {
		return models.Ticket{}, err
	}
	return t, nil
}
