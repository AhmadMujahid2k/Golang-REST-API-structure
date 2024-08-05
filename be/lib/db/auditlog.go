package db

import (
	psql "Golang-REST-API-structure/be/lib/psql"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
    ID         *string     `db:"id"`
    UserID     *string     `db:"u_id"`
    Action     *string     `db:"action"`
    CreatedAt  *time.Time  `db:"created_at"`
}

func (c *DbC) AuditLogs(
	tx *psql.Tx, 
	u_id uuid.UUID,
	action string,
) error {
    query := `
		INSERT INTO auditlog (
			u_id, 
			action, 
			created_at
		)	VALUES ($1, $2, $3)`
    
	err := psql.Exec(c.Pg, tx, query, u_id, action, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (c *DbC) GetAuditLogList(
    tx *psql.Tx, 
	offset int,
	limit int,
) ([]*AuditLog, error) {
    query := `
		SELECT * 
		FROM 
			auditlog
		ORDER BY 
			created_at DESC
		LIMIT $1 OFFSET $2`

	res , err := psql.Query[AuditLog](c.Pg, tx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	logs := make([]*AuditLog, len(*res))
	for i, log := range *res {
		logs[i] = &log
	}

	return logs, nil
}
