package main

import (
	"fmt"
	"github.com/Bayashat/TaskNinja/internal/data"
	"net/http"
	"time"
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

	// Create a new instance of the Movie struct, containing the ID we extracted from the URL and some dummy data.
	// Also notice that we deliberately haven't set a value for the UserID field.
	task := data.Task{
		ID:          id,
		CreatedAt:   time.Now(),
		Title:       "Golang Assignment-2",
		Description: "Create a project according to the book up to ch. 5. Database Setup and Configuration(Project should include 1-4 chapters).Send link to a git repository.Repositories with single commit will get -20%.",
		DueDate:     time.Date(2023, 10, 7, 23, 59, 0, 0, time.UTC),
		Priority:    "high",
		Status:      "in-process",
		Category:    "KBTU Tasks",
	}
	// Create an envelope{"task": task} instance and pass it to writeJSON(),
	// instead of passing the plain movie struct.
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
