package domain

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidTaskStatusTransition = errors.New("invalid task status transition")

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusBlocked    TaskStatus = "blocked"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityUrgent TaskPriority = "urgent"
)

type Task struct {
	ID          int64
	ProjectID   int64
	Title       string
	Description string
	Status      TaskStatus
	Priority    TaskPriority
	Source      Source
	Version     int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(
	projectID int64,
	title string,
	description string,
	priority TaskPriority,
	source Source,
) (*Task, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return nil, errors.New("title is required")
	}

	if projectID <= 0 {
		return nil, errors.New("projectID must be greater than 0")
	}

	if !ValidateTaskPriority(priority) {
		return nil, errors.New("invalid task priority")
	}

	if !ValidateSource(source) {
		return nil, errors.New("invalid source")
	}

	return &Task{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      TaskStatusTodo,
		Version:     1,
		Priority:    priority,
		Source:      source,
	}, nil
}

func (t *Task) ChangeStatus(to TaskStatus) error {
	if !CanTransitionTaskStatus(t.Status, to) {
		return ErrInvalidTaskStatusTransition
	}

	t.Status = to

	return nil
}

func ValidateTaskStatus(status TaskStatus) bool {
	return status == TaskStatusTodo ||
		status == TaskStatusBlocked ||
		status == TaskStatusCancelled ||
		status == TaskStatusDone ||
		status == TaskStatusInProgress
}

func ValidateTaskPriority(priority TaskPriority) bool {
	return priority == TaskPriorityHigh ||
		priority == TaskPriorityLow ||
		priority == TaskPriorityMedium ||
		priority == TaskPriorityUrgent
}

func ValidateSource(source Source) bool {
	return source == SourceClaude || source == SourceCodex || source == SourceHuman
}

func CanTransitionTaskStatus(from, to TaskStatus) bool {
	if !ValidateTaskStatus(from) || !ValidateTaskStatus(to) || from == to {
		return false
	}

	canFromTodo := from == TaskStatusTodo &&
		(to == TaskStatusInProgress || to == TaskStatusBlocked || to == TaskStatusCancelled || to == TaskStatusDone)

	canFromInProgress := from == TaskStatusInProgress &&
		(to == TaskStatusBlocked || to == TaskStatusCancelled || to == TaskStatusDone)

	canFromBlocked := from == TaskStatusBlocked &&
		(to == TaskStatusInProgress || to == TaskStatusCancelled)

	return canFromTodo || canFromInProgress || canFromBlocked
}
