package db

import (
	"Golang-REST-API-structure/be/lib/psql"
	"time"

	"github.com/google/uuid"
)

type SessionVals struct {
    Id   uuid.UUID
    Role string
}

type USignup struct {
    ID        uuid.UUID  `db:"id"`
    Fullname  string     `db:"full_name"`
    Email     string     `db:"email"`
    Password  string     `db:"password"`
    Role      string     `db:"role"`
}

type UPw struct {
    Password  string  `db:"password"`
}

func (c *DbC) Signup(tx *psql.Tx, u *USignup) error {
	query := ` 
		INSERT INTO u (
			id,
			full_name,
			email,
			password,
			role
		) VALUES ($1,$2,$3,$4,$5)`

	err := psql.Exec(c.Pg, tx, query, u.ID, u.Fullname, u.Email, u.Password, u.Role)
	if err != nil {
		return err
	}

	return nil
}

func (c *DbC) ULogin(tx *psql.Tx, email string) (*USignup, error) {
	query := `
		SELECT
			id,
			full_name,
			email,
			password,
			role 
		FROM 
			u
		WHERE 
			email = $1`
	 
	res, err := psql.QueryRow[USignup](c.Pg, tx, query, email)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *DbC) GetPwByID(tx *psql.Tx, userId uuid.UUID) (*UPw , error) {
	query := `
		SELECT 
			password 
		FROM 
			u 
		WHERE 
			id = $1`

	res, err := psql.QueryRow[UPw](c.Pg, tx, query, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *DbC) UpdatePwById(tx *psql.Tx, userId uuid.UUID, newPw string) error {
	query := `
		UPDATE 
			u
		SET 
			password = $1
		WHERE 
			id = $2`

	err := psql.Exec(c.Pg, tx, query, newPw, userId)
	if err != nil {
		return err
	}
	return nil
}

func (c *DbC) UpdateProfileById(tx *psql.Tx, userId uuid.UUID,phone string, dob time.Time, gender string, name string) error {
	query := `
		UPDATE 
			u
		SET 
			gender = $1, 
			dob = $2, 
			phone = $3, 
			full_name = $4
		WHERE 
			id = $5`

	err := psql.Exec(c.Pg, tx, query, gender, dob, phone, name, userId)
	if err != nil {
		return err
	}
	return nil
}
