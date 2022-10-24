// Filename: internal/data/todos.go

package data

import (
	"database/sql"
	"time"
	"errors"
	"fmt"
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

	args := []interface{}{
		todo.Item,
		todo.Description,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

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

		defer cancel()

		// Execute the query using QueryRow()
		err := m.DB.QueryRowContext(ctx, query, id).Scan(
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

	args := []interface{}{
		todo.Item,
		todo.Description,
		todo.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID)
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
	result, err := m.DB.ExecContext(ctx, query, id)
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
func (m TodoModel) GetAll(item string, description string, filters Filters) ([]*Todo, Metadata, error) {
	// Construct the query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, created_at, item, description
		FROM todolist
		WHERE (to_tsvector('simple', item) @@ plainto_tsquery('simple',$1) OR $1 = '')
		AND (to_tsvector('simple', description) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortOrder())

	// Create a 3-second-timout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query
	args := []interface{}{item, description, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Close the resultset
	defer rows.Close()
	totalRecords := 0
	// Initialize an empty slice to hold the School data
	lists := []*Todo{}
	// Iterate over the rows in the resultset
	for rows.Next() {
		var todo Todo
		// Scan the values from the row into school
		err := rows.Scan(
			&totalRecords,
			&todo.ID,
			&todo.CreatedAt,
			&todo.Item,
			&todo.Description,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the School to our slice
		lists = append(lists, &todo)
	}
	// Check for errors after looping through the resultset
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Return the slice of Schools
	return lists, metadata, nil
}