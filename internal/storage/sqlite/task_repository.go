package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

func (r *Repository) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var (
		taskID       int64
		createdAtRaw string
		updatedAtRaw string
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
		Scan(&taskID, &createdAtRaw, &updatedAtRaw)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	payload, err := marshalPayload(task.Status, 0, task.Version)
	if err != nil {
		return nil, err
	}

	err = insertActivityEvent(
		ctx,
		tx,
		task.ProjectID,
		domain.ActivityEventEntityTypeTask,
		taskID,
		"task_created",
		task.Source,
		payload,
	)
	if err != nil {
		return nil, err
	}

	createdAt, updatedAt, err := parseDates(createdAtRaw, updatedAtRaw)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create task: %w", err)
	}

	task.ID = taskID
	task.CreatedAt = createdAt
	task.UpdatedAt = updatedAt

	return task, nil
}

func (r *Repository) GetTaskByID(ctx context.Context, taskID int64) (*domain.Task, error) {
	query := `
		SELECT
			id,
			project_id,
			title,
			description,
			status,
			priority,
			source,
			version,
			created_at,
			updated_at
		FROM tasks
		WHERE id = $1
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, taskID)

	task, err := scanTaskRow(row.Scan, "get task", application.ErrTaskNotFound)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *Repository) UpdateTaskStatus(
	ctx context.Context,
	taskID int64,
	status domain.TaskStatus,
	expectedVersion int64,
	source domain.Source,
) (*domain.Task, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
	UPDATE tasks
	SET
			status = $1,
			version = version + 1,
			updated_at = datetime('now')
	WHERE id = $2
		AND version = $3
	RETURNING
			id,
			project_id,
			title,
			description,
			status,
			priority,
			source,
			version,
			created_at,
			updated_at
	`

	row := tx.QueryRowContext(ctx, query, status, taskID, expectedVersion)

	task, err := scanTaskRow(row.Scan, "update task", application.ErrTaskVersionConflict)
	if err != nil {
		return nil, err
	}

	payload, err := marshalPayload(task.Status, expectedVersion, task.Version)
	if err != nil {
		return nil, err
	}

	err = insertActivityEvent(
		ctx,
		tx,
		task.ProjectID,
		domain.ActivityEventEntityTypeTask,
		taskID,
		"task_status_updated",
		source,
		payload,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit update task: %w", err)
	}

	return task, nil
}

func (r *Repository) ListOpenTasksByProjectID(ctx context.Context, projectID int64) ([]domain.Task, error) {
	query := `
		SELECT
			id,
			project_id,
			title,
			description,
			status,
			priority,
			source,
			version,
			created_at,
			updated_at
		FROM tasks
		WHERE project_id = $1
			AND status IN ($2, $3, $4)
		ORDER BY created_at ASC, id ASC
	`
	rows, err := r.db.QueryContext(ctx, query, projectID, domain.TaskStatusTodo, domain.TaskStatusInProgress, domain.TaskStatusBlocked)
	if err != nil {
		return nil, fmt.Errorf("query open tasks by project ID %d: %w", projectID, err)
	}
	defer rows.Close()

	list := make([]domain.Task, 0)

	for rows.Next() {
		task, err := scanTaskRow(rows.Scan, "list open tasks", application.ErrTaskNotFound)
		if err != nil {
			return nil, err
		}

		list = append(list, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return list, nil
}

func scanTaskRow(scan func(dest ...any) error, reason string, noRowsError error) (*domain.Task, error) {
	var (
		task         domain.Task
		createdAtRaw string
		updatedAtRaw string
		statusRaw    string
		priorityRaw  string
		sourceRaw    string
	)

	err := scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&statusRaw,
		&priorityRaw,
		&sourceRaw,
		&task.Version,
		&createdAtRaw,
		&updatedAtRaw,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, noRowsError
	}

	if err != nil {
		return nil, fmt.Errorf("scan %s %d: %w", reason, task.ID, err)
	}

	createdAt, updatedAt, err := parseDates(createdAtRaw, updatedAtRaw)
	if err != nil {
		return nil, fmt.Errorf("%s ID %d: %w", reason, task.ID, err)
	}

	task.CreatedAt = createdAt
	task.UpdatedAt = updatedAt
	task.Status = domain.TaskStatus(statusRaw)
	task.Priority = domain.TaskPriority(priorityRaw)
	task.Source = domain.Source(sourceRaw)

	return &task, nil
}

func insertActivityEvent(
	ctx context.Context,
	tx *sql.Tx,
	projectID int64,
	entityType domain.ActivityEventEntityType,
	taskID int64,
	eventType string,
	source domain.Source,
	payloadJson []byte,
) error {
	query := `
	INSERT INTO activity_events (
		project_id,
		entity_type,
		entity_id,
		event_type,
		source,
		payload_json
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.ExecContext(ctx, query, projectID, entityType, taskID, eventType, source, string(payloadJson))
	if err != nil {
		return fmt.Errorf("create activity event: %w", err)
	}

	return nil
}

func marshalPayload(status domain.TaskStatus, previousVersion, version int64) ([]byte, error) {
	payload, err := json.Marshal(map[string]any{
		"status":           status,
		"previous_version": previousVersion,
		"version":          version,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal payload json: %w", err)
	}

	return payload, nil
}
