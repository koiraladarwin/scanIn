package postgres

import (
	"log"

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
	query := `INSERT INTO ticket (event_id, ticket_category_id, price, name,firebase_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, event_id, price, name;`
	var ticket models.Ticket
	err := p.sql.QueryRow(query, a.EventID, a.TicketCategoryID, a.Price, a.Name, a.FirebaseID).Scan(&ticket.ID, &ticket.TicketCategoryID, &ticket.EventID, &ticket.Price, &ticket.Name)
	if err != nil {
		return models.Ticket{}, err
	}
	return ticket, nil
}

func (p *PostgresDB) GetTickets(firebaseID string) ([]models.Ticket, error) {
	query := `SELECT id, event_id, ticket_category_id, price, name FROM ticket WHERE firebase_id = $1 AND deleted_at IS NULL;`
	rows, err := p.sql.Query(query, firebaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.EventID, &t.TicketCategoryID, &t.Price, &t.Name); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}
