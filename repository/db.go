package repository

import "database/sql"

// Queryer type
type Queryer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
}
