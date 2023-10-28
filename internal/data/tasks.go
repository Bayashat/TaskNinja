package data

import (
	"database/sql"
	"errors"
	"github.com/Bayashat/TaskNinja/internal/validator"
)

type Task struct {
	ID          int64      `json:"id"`          // Unique integer ID for the task
	CreatedAt   CustomTime `json:"created_at"`  // Timestamp for when the task is added to our database
	Title       string     `json:"title"`       // Task title
	Description string     `json:"description"` //  Task description
	//DueDate     CustomTime `json:"due_date"`    // Deadline or due date for the task
	Priority string `json:"priority"` // Task priority (e.g., high, medium, low)
	Status   string `json:"status"`   // Task status (e.g., to-do, in-progress, completed)
	Category string `json:"category"` // Task category or project it belongs to
	UserID   int64  `json:"user_id"`  // ID of the user who created the task (for multi-user support)
}

func ValidateTask(v *validator.Validator, task *Task) {
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(task.Description != "", "description", "must be provided")
	v.Check(len(task.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	//v.Check(!task.DueDate.IsZero(), "due_date", "must be provided")
	//v.Check(task.DueDate.Before(time.Date(2060, 1, 1, 0, 0, 0, 0, time.UTC)), "due_date", "must be before 2060")
	//v.Check(task.DueDate.After(time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC)), "due_date", "must be after 2023-10-07")
	v.Check(task.Priority != "", "priority", "must be provided")
	v.Check(task.Status != "", "status", "must be provided")
	v.Check(task.Category != "", "category", "must be provided")
}

// Define a TaskModel struct type which wraps a sql.DB connection pool.
type TaskModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the task table.
func (m TaskModel) Insert(task *Task) error {
	// Define the SQL query for inserting a new record in the task table and returning the system-generated data.
	query := `
		INSERT INTO tasks (title, description, priority, status, category)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, user_id`
	// Create an args slice containing the values for the placeholder parameters from the task struct.
	// Declaring this slice immediately next to our SQL query helps to make it nice
	// 		and clear *what values are being used where* in the query.
	args := []interface{}{task.Title, task.Description, task.Priority, task.Status, task.Category}
	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter
	// and scanning the system-generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&task.ID, &task.CreatedAt, &task.UserID)
}

// Add a placeholder method for fetching a specific record from the task table.
func (m TaskModel) Get(id int64) (*Task, error) {
	// The PostgreSQL bigserial type that we're using for the movie ID starts auto-incrementing at 1 by default,
	// so we know that no task will have ID values less than that.
	// To avoid making an unnecessary database call, we take a shortcut and return an ErrRecordNotFound error straight away.
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the task data.
	query := `
		SELECT id, created_at, title, description, priority, status, category, user_id
		FROM tasks
		WHERE id = $1`
	// Declare a Task struct to hold the data returned by the query.
	var task Task

	// Execute the query using the QueryRow() method, passing in the provided id value as a placeholder parameter,
	// and scan the response data into the fields of the Task struct.
	err := m.DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Status,
		&task.Category,
		&task.UserID,
	)
	// Handle any errors. If there was no matching task found, Scan() will return a sql.ErrNoRows error.
	// We check for this and return our custom ErrRecordNotFound error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &task, nil
}

// Add a placeholder method for updating a specific record in the task table.
func (m TaskModel) Update(task *Task) error {
	// Declare the SQL query for updating the record and returning the new version number.
	query := `
		UPDATE tasks
		SET title = $1, description = $2, priority = $3, status = $4, category = $5, user_id = $6
		WHERE id = $7
		RETURNING user_id`
	// Create an args slice containing the values for the placeholder parameters.
	args := []interface{}{
		task.Title,
		task.Description,
		task.Priority,
		task.Status,
		task.Category,
		task.UserID,
		task.ID,
	}
	// Use the QueryRow() method to execute the query,
	// passing in the args slice as a variadic parameter and scanning the new version value into the task struct.
	return m.DB.QueryRow(query, args...).Scan(&task.UserID)
}

// Add a placeholder method for deleting a specific record from the task table.
func (m TaskModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the task ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM tasks
		WHERE id = $1`
	// Execute the SQL query using the Exec() method, passing in the id variable as the value for the placeholder parameter.
	// The Exec() method returns a sql.Result object.
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected,
	//	we know that the tasks table didn't contain a record with the provided ID at the moment we tried to delete it.
	// In that case we return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
