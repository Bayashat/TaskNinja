package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the http.HandlerFunc() adapter,
	// and then set it as the custom error handler for 404 Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler
	// and set it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Add the route for the GET /v1/tasks endpoint.
	router.HandlerFunc(http.MethodGet, "/v1/tasks", app.listTasksHandler)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tasks", app.createTaskHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", app.showTaskHandler)
	// Require a PATCH request, rather than PUT.
	router.HandlerFunc(http.MethodPatch, "/v1/tasks/:id", app.updateTaskHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tasks/:id", app.deleteTaskHandler)

	return router
}
