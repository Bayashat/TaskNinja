package main

import (
	"fmt"
	"github.com/Bayashat/TaskNinja/internal/data"
	"github.com/Bayashat/TaskNinja/internal/validator"
	"net/http"
	"time"
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

	// Call the ValidateMovie() function and return a response containing the errors if any of the checks fail.
	if data.ValidateTask(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Call the Insert() method on our movies model, passing in a pointer to the validated movie struct.
	// This will create a record in the database and update the movie struct with the system-generated information.
	err = app.models.Tasks.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// When sending a HTTP response, we want to include a Location header to
	//		let the client know which URL they can find the newly-created resource at.
	// We make an empty http.Header map and then use the Set() method to add a new Location header,
	// 		interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tasks/%d", movie.ID))
	// Write a JSON response with a 201 Created status code, the movie data in the response body, and the Location header.
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
		//DueDate:     data.CustomTime(time.Date(2023, 10, 7, 23, 59, 0, 0, time.UTC)),
		Priority: "high",
		Status:   "in-process",
		Category: "KBTU Tasks",
	}
	// Create an envelope{"task": task} instance and pass it to writeJSON(),
	// instead of passing the plain movie struct.
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
