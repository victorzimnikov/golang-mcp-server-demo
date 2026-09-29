package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

func (r *Repository) ListProjects(ctx context.Context) ([]domain.Project, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM projects ORDER BY created_at ASC, id ASC"

	list := make([]domain.Project, 0)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query projects: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			project      domain.Project
			createdAtRaw string
			updatedAtRaw string
		)

		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &createdAtRaw, &updatedAtRaw); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}

		createdAt, updatedAt, err := parseDates(createdAtRaw, updatedAtRaw)
		if err != nil {
			return nil, fmt.Errorf("list projects: %w", err)
		}

		project.CreatedAt = createdAt
		project.UpdatedAt = updatedAt

		list = append(list, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}

	return list, nil
}

func (r *Repository) GetProjectByID(ctx context.Context, projectID int64) (*domain.Project, error) {
	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM projects
		WHERE id = $1
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, projectID)

	var (
		project      domain.Project
		createdAtRaw string
		updatedAtRaw string
	)

	err := row.Scan(&project.ID, &project.Name, &project.Description, &createdAtRaw, &updatedAtRaw)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, application.ErrProjectNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("get project by ID %d: %w", projectID, err)
	}

	createdAt, updatedAt, err := parseDates(createdAtRaw, updatedAtRaw)
	if err != nil {
		return nil, fmt.Errorf("get project by ID %d: %w", projectID, err)
	}

	project.CreatedAt = createdAt
	project.UpdatedAt = updatedAt

	return &project, nil
}
