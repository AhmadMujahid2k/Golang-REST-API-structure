package db

import (
	psql "Golang-REST-API-structure/be/lib/psql"
	"time"

	"github.com/google/uuid"
)

type Contact struct {
    ID             uuid.UUID  `db:"id"`
    UserID         uuid.UUID  `db:"u_id"`
    UserAgentID    uuid.UUID  `db:"u_agent_id"`
    AccountID      uuid.UUID  `db:"acc_id"`
    ContactNumber  string     `db:"contact_number"`
    CreatedAt      time.Time  `db:"created_at"`
}

func (c *DbC) DeleteContact(
	tx *psql.Tx, 
	contactId uuid.UUID,
) error {
    query := `
		DELETE 
		FROM 
			contact
		WHERE 
			id = $1`

    err := psql.Exec(c.Pg, tx, query, contactId)
    if err != nil {
        return err
    }

    return nil
}

func (c *DbC) GetContactList(
	tx *psql.Tx, 
	u_id uuid.UUID,
) ([]*Contact, error) {
    query := `
		SELECT * 
		FROM 
			contact
		WHERE 
			u_id = $1
		ORDER BY
			created_at DESC`

	res , err := psql.Query[Contact](c.Pg, tx, query, u_id)
	if err != nil {
		return nil, err
	}

	contacts := make([]*Contact, len(*res))
	for i, contact := range *res {
		contacts[i] = &contact
	}

	return contacts, nil
}

func (c *DbC) UploadContact(tx *psql.Tx, a *Contact) error {
    query := `
		INSERT INTO contact (
			id,
			u_id, 
			contact_number,
			created_at
		)	VALUES ($1, $2, $3, $4)`
    
	err := psql.Exec(c.Pg, tx, query, a.ID, a.UserID, a.ContactNumber, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (c *DbC) GetContact(
	tx *psql.Tx, 
	userID uuid.UUID,
	contactID uuid.UUID,
) (*Contact, error) {
    query := `
		SELECT * 
		FROM 
			contact
		WHERE 
			u_id = $1 AND 
			id = $2`

    contact, err := psql.QueryRow[Contact](c.Pg, tx, query, userID, contactID)
    if err != nil {
        return nil, err
    }

    return contact , nil
}
