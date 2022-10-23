// Filename: internal/data/schools.go

package data

import (
	"database/sql"

	"Quiz3.zioncastillo.net/internal/validator"

)

type Todo struct {
	ID           int64     `json:"id"`
	Item         string    `json:"name"`
	Category     string    `json:"level"`
}

func ValidateSchool(v *validator.Validator, todo *Todo) {
	// Use the Check() method to execute our validation checks

}

// Define a SchoolModel which wraps a sql.DB connection pool
type TodoModel struct {
	DB *sql.DB
}
