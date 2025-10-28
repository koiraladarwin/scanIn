package postgres

import (
	"encoding/json"
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
		e.name,
		e.description,
		e.start_time,
		e.end_time,
		e.event_organizer,
		e.location,
		COALESCE(
			json_agg(
				jsonb_build_object(
					'id', t.id,
					'ticket_category_id', t.ticket_category_id,
					'ticket_category_tag', tc.tag,
					'price', t.price,
					'attendees', COALESCE((
						SELECT json_agg(
							jsonb_build_object(
								'id', a.id,
								'attendee_id', u.id,
								'position', u.position,
								'company', u.company,
								'ticket_id', a.ticket_id,
								'deleted_at', a.deleted_at,
								'name', u.full_name,
								'email', u.gmail,
								'image_url', u.image_url
							)
						)
						FROM attendee a
						JOIN users u ON a.user_id = u.id
						WHERE a.ticket_id = t.id AND a.deleted_at IS NULL
					), '[]'::json)
				)
			) FILTER (WHERE t.id IS NOT NULL),
			'[]'::json
		) AS tickets
	FROM events e
	LEFT JOIN ticket t 
		ON e.id = t.event_id 
		AND t.deleted_at IS NULL 
		AND t.firebase_id = $1 
		AND t.price > 0
	LEFT JOIN ticket_category tc 
		ON t.ticket_category_id = tc.id
	WHERE e.deleted_at IS NULL 
		AND e.firebase_id = $1
	GROUP BY e.id
	ORDER BY e.start_time DESC;
	`

	rows, err := p.sql.Query(query, firebaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.EventsTicket

	for rows.Next() {
		var evt models.EventsTicket
		var ticketsJSON []byte

		if err := rows.Scan(
			&evt.Event.ID,
			&evt.Event.Name,
			&evt.Event.Description,
			&evt.Event.StartTime,
			&evt.Event.EndTime,
			&evt.Event.EventOrganizer,
			&evt.Event.Location,
			&ticketsJSON,
		); err != nil {
			return nil, err
		}

		// Unmarshal JSON into TicketDetails
		if err := json.Unmarshal(ticketsJSON, &evt.TicketDetails); err != nil {
			return nil, err
		}

		// Ensure Attendees slice is never null
		for i := range evt.TicketDetails {
			if evt.TicketDetails[i].Attendees == nil {
				evt.TicketDetails[i].Attendees = []models.AttendeeWithDetails{}
			}
		}

		results = append(results, evt)
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
		e.name,
		e.description,
		e.start_time,
		e.end_time,
		e.event_organizer,
		e.location,
		COALESCE(
			json_agg(
				jsonb_build_object(
					'id', t.id,
					'ticket_category_id', t.ticket_category_id,
					'ticket_category_tag', tc.tag,
					'price', t.price,
					'attendees', COALESCE((
						SELECT json_agg(
							jsonb_build_object(
								'id', a.id,
								'attendee_id', u.id,
								'position', u.position,
								'company', u.company,
								'ticket_id', a.ticket_id,
								'deleted_at', a.deleted_at,
								'name', u.full_name,
								'email', u.gmail,
								'image_url', u.image_url
							)
						)
						FROM attendee a
						JOIN users u ON a.user_id = u.id
						WHERE a.ticket_id = t.id AND a.deleted_at IS NULL
					), '[]'::json)
				)
			) FILTER (WHERE t.id IS NOT NULL),
			'[]'::json
		) AS tickets
	FROM events e
	LEFT JOIN ticket t 
		ON e.id = t.event_id 
		AND t.deleted_at IS NULL 
		AND t.firebase_id = $1 
		AND t.price > 0
	LEFT JOIN ticket_category tc 
		ON t.ticket_category_id = tc.id
	WHERE e.deleted_at IS NULL 
		AND e.firebase_id = $1
	GROUP BY e.id
	ORDER BY e.start_time DESC;
	`

	rows, err := p.sql.Query(query, firebaseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.EventsInvitee

	for rows.Next() {
		var evt models.EventsInvitee
		var ticketsJSON []byte

		if err := rows.Scan(
			&evt.Event.ID,
			&evt.Event.Name,
			&evt.Event.Description,
			&evt.Event.StartTime,
			&evt.Event.EndTime,
			&evt.Event.EventOrganizer,
			&evt.Event.Location,
			&ticketsJSON,
		); err != nil {
			return nil, err
		}

		// Unmarshal JSON into TicketDetails
		if err := json.Unmarshal(ticketsJSON, &evt.InviteeDetails); err != nil {
			return nil, err
		}

		// Ensure Attendees slice is never null
		for i := range evt.InviteeDetails {
			if evt.InviteeDetails[i].Attendees == nil {
				evt.InviteeDetails[i].Attendees = []models.AttendeeWithDetails{}
			}
		}

		results = append(results, evt)
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
