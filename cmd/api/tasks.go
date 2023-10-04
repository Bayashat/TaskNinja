package main

import (
	"fmt"
	"net/http"
)

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Creating a new task...")
}

// Add a showTaskHandler for the "GET /v1/task/:id" endpoint.
// For now, we retrieve the interpolated "id" parameter from the current URL and include it in a placeholder response.

func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "show the details of task: %d\n", id)
}
