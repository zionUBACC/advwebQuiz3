// Filename: cmd/api/helpers.go

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

//Define a new type name envelope
type envelope map[string]interface{}

func (app *application) readIDParam(r *http.Request) (int64, error){
	//use the "paramsfromcontext()" function to get the request context a slice
	params := httprouter.ParamsFromContext(r.Context())
	//GET the value of the "id" parameter
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		
		return 0, errors.New("Invalid Id Parameter")
	}
	return id, nil
}
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	//Convert our map into a JSON object
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	// Add a newline to make viewing on the terminal easier
	js = append(js, '\n')
	// Add the headers
	for key, value := range headers {
		w.Header()[key] = value
	}
	// Specify that we will serve our responses using JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Write the []byte slice containing the JSON response body
	w.Write(js)
	return nil
}
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	
	//use http.maxbytesreader() to limit the size of the request body to 
	//1 MB 2^20
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()	
	
	//decode the request body into the destination
	err := dec.Decode(dst)
	
	//check for a bad request
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		
		// switch to check the errror
		switch {
			//check for syntax errors
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains formed JSON(at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		//check for wrong types passed by the client
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		//empty body
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		//unmappable fields
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		//too large file
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		//pass a non-nil pointer error
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		//default
		default:
			return err
		}
	}
		//call decode again
	err = dec.Decode(&struct {}{})
	if err != io.EOF {
		return errors.New("body must only contain a single value")
	}
		

	
	return nil
}