package main

import (
	"encoding/json"
	"fmt"
	"github.com/Bayashat/TaskNinja/internal/data"
	"net/http"
	"time"
)

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the HTTP request body
	// (note that the field names and types in the struct are a subset of the Movie struct that we created earlier).
	// This struct will be our *target  decode destination*.
	var input struct {
		Title       string          `json:"title"`
		Description string          `json:"description"`
		DueDate     data.CustomTime `json:"due_date"`
		Priority    string          `json:"priority"`
		Status      string          `json:"status"`
		Category    string          `json:"category"`
	}
	// Initialize a new json.Decoder instance which reads from the request body,
	// 		and then use the Decode() method to decode the body contents into the input struct.
	// Importantly, notice that when we call Decode() we pass a *pointer* to the input struct as the target decode destination.
	// If there was an error during decoding,
	// 	we also use our generic errorResponse() helper to send the client a 400 Bad Request response containing the error message.
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	// Dump the contents of the input struct in a HTTP response.
	fmt.Fprintf(w, "%+v\n", input)
}

// Add a showTaskHandler for the "GET /v1/task/:id" endpoint.
// For now, we retrieve the interpolated "id" parameter from the current URL and include it in a placeholder response.

func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		// Use the new notFoundResponse() helper.
		app.notFoundResponse(w, r)
		return
	}

	// Create a new instance of the Movie struct, containing the ID we extracted from the URL and some dummy data.
	// Also notice that we deliberately haven't set a value for the UserID field.
	task := data.Task{
		ID:          id,
		CreatedAt:   data.CustomTime(time.Now()),
		Title:       "Golang Assignment-2",
		Description: "Create a project according to the book up to ch. 5. Database Setup and Configuration(Project should include 1-4 chapters).Send link to a git repository.Repositories with single commit will get -20%.",
		DueDate:     data.CustomTime(time.Date(2023, 10, 7, 23, 59, 0, 0, time.UTC)),
		Priority:    "high",
		Status:      "in-process",
		Category:    "KBTU Tasks",
	}
	// Create an envelope{"task": task} instance and pass it to writeJSON(),
	// instead of passing the plain movie struct.
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
