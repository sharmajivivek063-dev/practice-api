package models

import (
	"time"
)

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=150"`
	Description string `json:"description"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high"`
	DueDate     string `json:"due_date"`
}

type TaskResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	UserID      int64  `json:"user_id"`
}

type TaskListItem struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	DueDate  string `json:"due_date"`
}

type TaskDetailResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     string    `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ********************************************************************
type UpdateTaskRequest struct {
	Title       *string `json:"title" binding:"required"`
	Description *string `json:"description" binding:"required"`
	Status      *string `json:"status" binding:"required"`
	Priority    *string `json:"priority" binding:"required"`
	DueDate     *string `json:"due_date" binding:"required"`
}

type UpdateTaskResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
}

// ***********************************************************************

type TaskStatsResponse struct {
	TotalTasks   int `json:"total_tasks"`
	Pending      int `json:"pending"`
	InProgress   int `json:"in_progress"`
	Completed    int `json:"completed"`
	HighPriority int `json:"high_priority"`
	Overdue      int `json:"overdue"`
}
