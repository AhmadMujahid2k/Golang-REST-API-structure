package db

import psql "Golang-REST-API-structure/be/lib/psql"

type DbC struct {
	Pg *psql.Postgres
}

func Init(
	Pg *psql.Postgres,
) *DbC {
	return &DbC{
		Pg: Pg,
	}
}
