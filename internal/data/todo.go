// Filename: internal/data/todos.go

package data

import (
	"database/sql"
	"time"
	"errors"
	"context"

	"Quiz3.zioncastillo.net/internal/validator"

)

type Todo struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"-"`
	Item         string    `json:"item"`
	Description  string    `json:"description"`
}

func ValidateItem(v *validator.Validator, todo *Todo) {
	// Use the Check() method to execute our validation checks
	v.Check(todo.Item != "", "name", "must be provided")
	v.Check(len(todo.Item) <= 200, "Item", "must not be more than 200 bytes long")
	v.Check(len(todo.Description) <= 2000, "level", "must not be more than 2000 bytes long")
}

// Define a TodoModel which wraps a sql.DB connection pool
type TodoModel struct {
	DB *sql.DB
}

func (m TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todolist (item, description)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{
		todo.Item,
		todo.Description,
	}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt)
}

func (m TodoModel) Get(id int64) (*Todo, error) {
		// Ensure that there is a valid id
		if id < 1 {
			return nil, ErrRecordNotFound
		}
		// Create the query
		query := `
			SELECT id, created_at, item, description
			FROM todolist
			WHERE id = $1
		`
		// Declare a School variable to hold the returned data
		var todo Todo

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		defer cancel

		// Execute the query using QueryRow()
		err := m.DB.QueryRowContext(query, id).Scan(
			&todo.ID,
			&todo.CreatedAt,
			&todo.Item,
			&todo.Description,
		)
		// Handle any errors
		if err != nil {
			// Check the type of error
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return nil, ErrRecordNotFound
			default:
				return nil, err
			}
		}
		// Success
		
	return &todo, nil
	
}

// Update() allows us to edit/alter a specific Todo
func (m TodoModel) Update(todo *Todo) error {
		// Create a query
		query := `
		UPDATE todolist
		SET item = $1, description = $2
		WHERE id = $3
		RETURNING id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	
	defer cancel()

	args := []interface{}{
		todo.Item,
		todo.Description,
		todo.ID,
	}
	return m.DB.QueryRowContext(query, args...).Scan(&todo.ID)
}

// Delete() removes a specific Todo
func (m TodoModel) Delete(id int64) error {
	// Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	// Create the delete query
	query := `
		DELETE FROM todolist
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// Execute the query
	result, err := m.DB.ExecContext(query, id)
	if err != nil {
		return err
	}
	// Check how many rows were affected by the delete operation. We
	// call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}