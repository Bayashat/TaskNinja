package data

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
