// Filename: cmd/api/routes.go

package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)
func (app *application) routes () *httprouter.Router{
	
	router := httprouter.New()
	
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/list", app.createTodoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/list", app.listTodoListHandler)
	router.HandlerFunc(http.MethodGet, "/v1/list/:id", app.showTodoHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/list/:id", app.updateTodoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/list/:id", app.deleteTodoHandler)

	return router
}