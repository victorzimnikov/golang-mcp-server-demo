package sqlite

import (
	"fmt"
	"time"
)

func parseDates(createdAtRaw, updatedAtRaw string) (time.Time, time.Time, error) {
	createdAt, err := parseTime(createdAtRaw)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse created_at: %w", err)
	}

	updatedAt, err := parseTime(updatedAtRaw)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse updated_at: %w", err)
	}

	return createdAt, updatedAt, nil
}
