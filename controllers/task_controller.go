package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"practice-api/connection"
	"practice-api/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {

	// Step 1: Get authenticated user ID from JWT middleware

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized access",
		})
		return
	}

	userID, ok := userIDValue.(int64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Step 2: Bind and validate JSON request

	var req models.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Step 3: Validate due date if provided

	var dueDate interface{}

	if req.DueDate != "" {

		parsedDate, err := time.Parse("2006-01-02", req.DueDate)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid due date format. Use YYYY-MM-DD",
			})
			return
		}

		dueDate = parsedDate
	}

	// Step 4: Insert task into database
	// Status is automatically set to 'pending' by MySQL.

	query := `
        INSERT INTO tasks
        (title, description, priority, due_date, user_id)
        VALUES (?, ?, ?, ?, ?)
    `

	result, err := connection.DB.Exec(
		query,
		req.Title,
		req.Description,
		req.Priority,
		dueDate,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create task",
		})
		return
	}

	// Step 5: Get newly created task ID

	taskID, err := result.LastInsertId()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve task ID",
		})
		return
	}

	// Step 6: Fetch the created task from database

	var task models.TaskResponse

	selectQuery := `
        SELECT
            id,
            title,
            description,
            status,
            priority,
            COALESCE(DATE_FORMAT(due_date, '%Y-%m-%d'), ''),
            user_id
        FROM tasks
        WHERE id = ? AND user_id = ?
    `

	err = connection.DB.QueryRow(
		selectQuery,
		taskID,
		userID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.UserID,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Created task not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch created task",
		})
		return
	}

	// Step 7: Return success response

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Task created successfully",
		"data":    task,
	})
}

// ****************************************************************
func ListTasks(c *gin.Context) {

	// Step 1: Get user ID from JWT middleware

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized access",
		})
		return
	}

	userID, ok := userIDValue.(int64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Step 2: Read pagination parameters

	page := 1
	limit := 10

	if pageValue := c.Query("page"); pageValue != "" {
		parsedPage, err := strconv.Atoi(pageValue)

		if err != nil || parsedPage < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Page must be a positive integer",
			})
			return
		}

		page = parsedPage
	}

	if limitValue := c.Query("limit"); limitValue != "" {
		parsedLimit, err := strconv.Atoi(limitValue)

		if err != nil || parsedLimit < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Limit must be a positive integer",
			})
			return
		}

		limit = parsedLimit

		// Maximum limit is 100
		if limit > 100 {
			limit = 100
		}
	}

	// Step 3: Read filters

	status := strings.TrimSpace(c.Query("status"))
	priority := strings.TrimSpace(c.Query("priority"))
	search := strings.TrimSpace(c.Query("search"))

	// Validate priority if provided

	if priority != "" &&
		priority != "low" &&
		priority != "medium" &&
		priority != "high" {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Priority must be low, medium, or high",
		})
		return
	}

	// Step 4: Build WHERE conditions

	conditions := []string{"user_id = ?"}
	args := []interface{}{userID}

	if status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	if priority != "" {
		conditions = append(conditions, "priority = ?")
		args = append(args, priority)
	}

	if search != "" {
		conditions = append(
			conditions,
			"(title LIKE ? OR description LIKE ?)",
		)

		searchValue := "%" + search + "%"

		args = append(args, searchValue, searchValue)
	}

	whereClause := strings.Join(conditions, " AND ")

	// Step 5: Count records after applying filters

	countQuery := `
        SELECT COUNT(*)
        FROM tasks
        WHERE ` + whereClause

	var totalRecords int

	err := connection.DB.QueryRow(
		countQuery,
		args...,
	).Scan(&totalRecords)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to count tasks",
		})
		return
	}

	// Step 6: Calculate pagination

	totalPages := (totalRecords + limit - 1) / limit

	offset := (page - 1) * limit

	// Step 7: Fetch tasks

	listQuery := `
        SELECT
            id,
            title,
            status,
            priority,
            COALESCE(DATE_FORMAT(due_date, '%Y-%m-%d'), '')
        FROM tasks
        WHERE ` + whereClause + `
        ORDER BY created_at DESC, id DESC
        LIMIT ? OFFSET ?
    `

	listArgs := append([]interface{}{}, args...)
	listArgs = append(listArgs, limit, offset)

	rows, err := connection.DB.Query(
		listQuery,
		listArgs...,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch tasks",
		})
		return
	}

	defer rows.Close()

	// Step 8: Store task records

	tasks := []models.TaskListItem{}

	for rows.Next() {

		var task models.TaskListItem

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Status,
			&task.Priority,
			&task.DueDate,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to read task data",
			})
			return
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error while reading tasks",
		})
		return
	}

	// Step 9: Return response

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tasks,
		"pagination": gin.H{
			"page":          page,
			"limit":         limit,
			"total_records": totalRecords,
			"total_pages":   totalPages,
		},
	})
}

// ***********************************************************************************

func GetTask(c *gin.Context) {

	// Step 1: Get authenticated user ID from JWT

	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized access",
		})
		return
	}

	userID, ok := userIDValue.(int64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Step 2: Get task ID from URL parameter

	taskID, err := strconv.Atoi(c.Param("id"))

	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid task ID",
		})
		return
	}

	// Step 3: Fetch task belonging to authenticated user

	var task models.TaskDetailResponse

	query := `
        SELECT
            id,
            title,
            description,
            status,
            priority,
            COALESCE(DATE_FORMAT(due_date, '%Y-%m-%d'), ''),
            created_at,
            updated_at
        FROM tasks
        WHERE id = ? AND user_id = ?
    `

	err = connection.DB.QueryRow(
		query,
		taskID,
		userID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	// Step 4: Handle task not found

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Task not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch task",
		})
		return
	}

	// Step 5: Return task details

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    task,
	})
}

// ************************************************************************************

func UpdateTask(c *gin.Context) {
	// 1. Get task ID from URL
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid task ID",
		})
		return
	}

	// 2. Get authenticated user ID from JWT middleware
	userValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID, ok := userValue.(int64)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// 3. Bind JSON request
	var req models.UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "All five fields are required",
			"error":   err.Error(),
		})
		return
	}

	// 4. Validate title
	if strings.TrimSpace(*req.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Title cannot be empty",
		})
		return
	}

	// 5. Validate status
	validStatuses := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"completed":   true,
	}

	if !validStatuses[*req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Status must be pending, in_progress, or completed",
		})
		return
	}

	// 6. Validate priority
	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}

	if !validPriorities[*req.Priority] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Priority must be low, medium, or high",
		})
		return
	}

	// 7. Validate due_date
	dueDate, err := time.Parse("2006-01-02", *req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "due_date must be in YYYY-MM-DD format",
		})
		return
	}

	// 8. Update only the task owned by this user
	result, err := connection.DB.Exec(`
		UPDATE tasks
		SET title = ?,
		    description = ?,
		    status = ?,
		    priority = ?,
		    due_date = ?
		WHERE id = ? AND user_id = ?
	`,
		strings.TrimSpace(*req.Title),
		*req.Description,
		*req.Status,
		*req.Priority,
		dueDate.Format("2006-01-02"),
		taskID,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update task",
			"error":   err.Error(),
		})
		return
	}

	// 9. Check whether a row was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Could not verify task update",
		})
		return
	}

	// MySQL may report zero affected rows if the values were unchanged.
	// Check whether the task exists and belongs to this user.
	if rowsAffected == 0 {
		var existingID int

		err := connection.DB.QueryRow(
			"SELECT id FROM tasks WHERE id = ? AND user_id = ?",
			taskID,
			userID,
		).Scan(&existingID)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Task not found",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to verify task",
			})
			return
		}
	}

	// 10. Return updated task
	updatedTask := models.UpdateTaskResponse{
		ID:          taskID,
		Title:       strings.TrimSpace(*req.Title),
		Description: *req.Description,
		Status:      *req.Status,
		Priority:    *req.Priority,
		DueDate:     dueDate.Format("2006-01-02"),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task updated successfully",
		"data":    updatedTask,
	})
}

// *********************************************************************************

func DeleteTask(c *gin.Context) {
	// 1. Get task ID from URL
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid task ID",
		})
		return
	}

	// 2. Get authenticated user ID from JWT middleware
	userValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID, ok := userValue.(int64)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// 3. Delete only the task owned by this user
	result, err := connection.DB.Exec(
		"DELETE FROM tasks WHERE id = ? AND user_id = ?",
		taskID,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete task",
		})
		return
	}

	// 4. Check affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Could not verify task deletion",
		})
		return
	}

	// 5. Task not found or not owned by this user
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Task not found",
		})
		return
	}

	// 6. Successful deletion
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task deleted successfully",
	})
}

// ************************************************************************************

func GetTaskStats(c *gin.Context) {
	// 1. Get authenticated user ID from JWT middleware
	userValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID, ok := userValue.(int64)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// 2. Prepare response
	var stats models.TaskStatsResponse

	// 3. Execute aggregate SQL query
	query := `
    SELECT
        COUNT(*),

        COALESCE(SUM(CASE
            WHEN status = 'pending' THEN 1
            ELSE 0
        END), 0),

        COALESCE(SUM(CASE
            WHEN status = 'in_progress' THEN 1
            ELSE 0
        END), 0),

        COALESCE(SUM(CASE
            WHEN status = 'completed' THEN 1
            ELSE 0
        END), 0),

        COALESCE(SUM(CASE
            WHEN priority = 'high' THEN 1
            ELSE 0
        END), 0),

        COALESCE(SUM(CASE
            WHEN due_date < CURDATE()
            AND status != 'completed' THEN 1
            ELSE 0
        END), 0)

    FROM tasks
    WHERE user_id = ?
`

	err := connection.DB.QueryRow(query, userID).Scan(
		&stats.TotalTasks,
		&stats.Pending,
		&stats.InProgress,
		&stats.Completed,
		&stats.HighPriority,
		&stats.Overdue,
	)
	if err != nil {
		log.Println("GetTaskStats database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch task statistics",
		})
		return
	}

	log.Println("Task stats:", stats)

	// 4. Return statistics
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})

}
