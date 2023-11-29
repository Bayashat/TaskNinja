package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Update the routes() method to return a http.Handler instead of a *httprouter.Router.
func (app *application) routes() http.Handler {
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

	// Add the route for the POST /v1/users endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	// Add the route for the PUT /v1/users/activated endpoint.
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	// Wrap the router with the rateLimit() middleware.
	return app.recoverPanic(app.rateLimit(router))
}
