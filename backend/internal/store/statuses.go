package store

import (
	"context"
	"database/sql"

	"jobtracker/internal/models"
)

type StatusStore struct {
	db *sql.DB
}

func NewStatusStore(db *sql.DB) *StatusStore {
	return &StatusStore{db: db}
}

func (s *StatusStore) List(ctx context.Context) ([]models.ApplicationStatus, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, code, label, sort_order FROM application_statuses ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.ApplicationStatus
	for rows.Next() {
		var st models.ApplicationStatus
		if err := rows.Scan(&st.ID, &st.Code, &st.Label, &st.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, rows.Err()
}
