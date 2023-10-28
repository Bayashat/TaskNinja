package main

import (
	"errors"
	"fmt"
	"github.com/Bayashat/TaskNinja/internal/data"
	"github.com/Bayashat/TaskNinja/internal/validator"
	"net/http"
)

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the HTTP request body
	// (note that the field names and types in the struct are a subset of the Movie struct that we created earlier).
	// This struct will be our *target  decode destination*.
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		//DueDate     data.CustomTime `json:"due_date"`
		Priority string `json:"priority"`
		Status   string `json:"status"`
		Category string `json:"category"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the input struct to a new Movie struct.
	movie := &data.Task{
		Title:       input.Title,
		Description: input.Description,
		//DueDate:     input.DueDate,
		Priority: input.Priority,
		Status:   input.Status,
		Category: input.Category,
	}

	// Initialize a new Validator.
	v := validator.New()

	// Call the ValidateTask() function and return a response containing the errors if any of the checks fail.
	if data.ValidateTask(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Call the Insert() method on our tasks model, passing in a pointer to the validated task struct.
	// This will create a record in the database and update the task struct with the system-generated information.
	err = app.models.Tasks.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// When sending a HTTP response, we want to include a Location header to
	//		let the client know which URL they can find the newly-created resource at.
	// We make an empty http.Header map and then use the Set() method to add a new Location header,
	// 		interpolating the system-generated ID for our new task in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tasks/%d", movie.ID))
	// Write a JSON response with a 201 Created status code, the task data in the response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"task": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Add a showTaskHandler for the "GET /v1/task/:id" endpoint.
// For now, we retrieve the interpolated "id" parameter from the current URL and include it in a placeholder response.

func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Call the Get() method to fetch the data for a specific task.
	// We also need to use the errors.Is() function to check if it returns a data.ErrRecordNotFound error,
	// in which case we send a 404 Not Found response to the client.
	task, err := app.models.Tasks.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
