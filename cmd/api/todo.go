// Filename: cms/api/todo.go
package main
import (
	"fmt"
	"errors"
	"net/http"

	"Quiz3.zioncastillo.net/internal/data"
	"Quiz3.zioncastillo.net/internal/validator"
)
// createSchoolHandler for the "POST /v1/schools" endpoint
func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Our target decode destination
	var input struct{
		Item        string   `json:"item"`
		Descript    string   `json:"description"`
	}
	// Initialize a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	 //Copy the values from the input struct to a new School struct
	 todo := &data.Todo{
	 	Item: input.Item,
	 	Description: input.Descript,
	}
	// Initialize a new Validator instance
	v := validator.New()
	// Check the map to determine if there were any validation errors
	if data.ValidateItem(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// // Display the request
	// fmt.Fprintf(w, "%+v\n", input)
	// Create a School
	err = app.models.Todo.Insert(todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// Create a Location header for the newly created resource/School
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/todo/%d", todo.ID))
	// Write the JSON response with 201 - Created status code with the body
	// being the School data and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"item": todo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showSchoolHandler for the "GET /v1/schools/:id" endpoint
func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	todo, err := app.models.Todo.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"item": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// This method does a complete replacement
	// Get the id for the school that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the orginal record from the database
	todo, err := app.models.Todo.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Create an input struct to hold data read in fro mteh client
	var input struct {
		Item       string   `json:"item"`
		Descript   string   `json:"description"`
	}

	// Initialize a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy / Update the fields / values in the school variable using the fields
	// in the input struct
	todo.Item = input.Item
	todo.Description = input.Descript

	// Perform validation on the updated School. If validation fails, then
	// we send a 422 - Unprocessable Entity respose to the client
	// Initialize a new Validator instance
	v := validator.New()

	// Check the map to determine if there were any validation errors
	if data.ValidateItem(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updated School record to the Update() method
	err = app.models.Todo.Update(todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id for the school that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the School from the database. Send a 404 Not Found status code to the
	// client if there is no matching record
	err = app.models.Todo.Delete(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return 200 Status OK to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Item successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}