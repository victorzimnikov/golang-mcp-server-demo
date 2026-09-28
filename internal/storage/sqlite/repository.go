package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

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

		project.CreatedAt, err = parseTime(createdAtRaw)
		if err != nil {
			return nil, fmt.Errorf("parse project createAt: %w", err)
		}

		project.UpdatedAt, err = parseTime(updatedAtRaw)
		if err != nil {
			return nil, fmt.Errorf("parse project updatedAt: %w", err)
		}

		list = append(list, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}

	return list, nil
}

func parseTime(value string) (time.Time, error) {
	return time.ParseInLocation(
		time.DateTime,
		value,
		time.UTC,
	)
}
