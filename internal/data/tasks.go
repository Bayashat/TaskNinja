package data

import (
	"database/sql"
	"github.com/Bayashat/TaskNinja/internal/validator"
	"time"
)

type Task struct {
	ID          int64      `json:"id"`          // Unique integer ID for the task
	CreatedAt   CustomTime `json:"created_at"`  // Timestamp for when the task is added to our database
	Title       string     `json:"title"`       // Task title
	Description string     `json:"description"` //  Task description
	DueDate     CustomTime `json:"due_date"`    // Deadline or due date for the task
	Priority    string     `json:"priority"`    // Task priority (e.g., high, medium, low)
	Status      string     `json:"status"`      // Task status (e.g., to-do, in-progress, completed)
	Category    string     `json:"category"`    // Task category or project it belongs to
	UserID      int64      `json:"user_id"`     // ID of the user who created the task (for multi-user support)
}

func ValidateTask(v *validator.Validator, task *Task) {
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(task.Description != "", "description", "must be provided")
	v.Check(len(task.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	v.Check(!task.DueDate.IsZero(), "due_date", "must be provided")
	v.Check(task.DueDate.Before(time.Date(2060, 1, 1, 0, 0, 0, 0, time.UTC)), "due_date", "must be before 2060")
	v.Check(task.DueDate.After(time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC)), "due_date", "must be after 2023-10-07")
	v.Check(task.Priority != "", "priority", "must be provided")
	v.Check(task.Status != "", "status", "must be provided")
	v.Check(task.Category != "", "category", "must be provided")
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type TaskModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m TaskModel) Insert(task *Task) error {
	return nil
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m TaskModel) Get(id int64) (*Task, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m TaskModel) Update(task *Task) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m TaskModel) Delete(id int64) error {
	return nil
}
