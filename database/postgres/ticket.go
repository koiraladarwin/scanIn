package postgres

import "github.com/koiraladarwin/scanin/models"

func (p *PostgresDB) CreateTicket(a models.TicketRequest) (models.Ticket, error) {
	query := `INSERT INTO ticket (event_id, price, name,firebase_id) VALUES ($1, $2, $3, $4) RETURNING id, event_id, price, name;`
	var ticket models.Ticket
	err := p.sql.QueryRow(query, a.EventID, a.Price, a.Name, a.FirebaseID).Scan(&ticket.ID, &ticket.EventID, &ticket.Price, &ticket.Name)
	if err != nil {
		return models.Ticket{}, err
	}
	return ticket, nil
}
