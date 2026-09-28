package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

func (r *Repository) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var (
		taskID    int64
		createdAt string
		updatedAt string
	)

	query := `
	INSERT INTO tasks (
		project_id,
    title,
    description,
    status,
    priority,
    source,
    version
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, updated_at`
	err = tx.
		QueryRowContext(ctx, query, task.ProjectID, task.Title, task.Description, task.Status, task.Priority, task.Source, task.Version).
		Scan(&taskID, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	query = `
	INSERT INTO activity_events (
		project_id,
		entity_type,
		entity_id,
		event_type,
		source
	)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.ExecContext(ctx, query, task.ProjectID, domain.ActivityEventEntityTypeTask, taskID, "task_created", task.Source)
	if err != nil {
		return nil, fmt.Errorf("create activity event: %w", err)
	}

	parsedCreatedAt, err := parseTime(createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse createdAt: %w", err)
	}
	parsedUpdatedAt, err := parseTime(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse updatedAt: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create task: %w", err)
	}

	task.ID = taskID
	task.CreatedAt = parsedCreatedAt
	task.UpdatedAt = parsedUpdatedAt

	return task, nil
}
