package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"jobtracker/internal/models"
)

type ApplicationStore struct {
	db *sql.DB
}

func NewApplicationStore(db *sql.DB) *ApplicationStore {
	return &ApplicationStore{db: db}
}

type ApplicationFilter struct {
	StatusID *uint8
	Company  string
}

func (s *ApplicationStore) List(ctx context.Context, f ApplicationFilter) ([]models.Application, error) {
	query := `SELECT id, company_name, position_title, job_description, job_posting_url, location,
		work_mode, salary_min, salary_max, source, current_status_id, priority, applied_date,
		notes, created_at, updated_at, archived_at
		FROM applications WHERE archived_at IS NULL`
	var args []any

	if f.StatusID != nil {
		query += " AND current_status_id = ?"
		args = append(args, *f.StatusID)
	}
	if f.Company != "" {
		query += " AND company_name LIKE ?"
		args = append(args, "%"+f.Company+"%")
	}
	query += " ORDER BY created_at DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Application
	for rows.Next() {
		app, err := scanApplication(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, app)
	}
	return out, rows.Err()
}

func (s *ApplicationStore) Get(ctx context.Context, id uint64) (models.Application, error) {
	query := `SELECT id, company_name, position_title, job_description, job_posting_url, location,
		work_mode, salary_min, salary_max, source, current_status_id, priority, applied_date,
		notes, created_at, updated_at, archived_at
		FROM applications WHERE id = ?`
	row := s.db.QueryRowContext(ctx, query, id)
	app, err := scanApplication(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Application{}, ErrNotFound
	}
	return app, err
}

type CreateApplicationInput struct {
	CompanyName     string
	PositionTitle   string
	JobDescription  *string
	JobPostingURL   *string
	Location        *string
	WorkMode        *string
	SalaryMin       *uint32
	SalaryMax       *uint32
	Source          *string
	CurrentStatusID uint8
	Priority        string
	AppliedDate     *string
	Notes           *string
}

func (s *ApplicationStore) Create(ctx context.Context, in CreateApplicationInput) (uint64, error) {
	query := `INSERT INTO applications
		(company_name, position_title, job_description, job_posting_url, location, work_mode,
		 salary_min, salary_max, source, current_status_id, priority, applied_date, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := s.db.ExecContext(ctx, query,
		in.CompanyName, in.PositionTitle, in.JobDescription, in.JobPostingURL, in.Location, in.WorkMode,
		in.SalaryMin, in.SalaryMax, in.Source, in.CurrentStatusID, in.Priority, in.AppliedDate, in.Notes,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

type UpdateApplicationInput struct {
	CompanyName    string
	PositionTitle  string
	JobDescription *string
	JobPostingURL  *string
	Location       *string
	WorkMode       *string
	SalaryMin      *uint32
	SalaryMax      *uint32
	Source         *string
	Priority       string
	AppliedDate    *string
	Notes          *string
}

// Update changes everything about an application EXCEPT its status.
// Status changes go through UpdateStatus so application_status_history stays
// in sync - there is deliberately no other path that touches
// current_status_id.
func (s *ApplicationStore) Update(ctx context.Context, id uint64, in UpdateApplicationInput) error {
	query := `UPDATE applications SET
		company_name = ?, position_title = ?, job_description = ?, job_posting_url = ?,
		location = ?, work_mode = ?, salary_min = ?, salary_max = ?, source = ?,
		priority = ?, applied_date = ?, notes = ?
		WHERE id = ?`

	res, err := s.db.ExecContext(ctx, query,
		in.CompanyName, in.PositionTitle, in.JobDescription, in.JobPostingURL,
		in.Location, in.WorkMode, in.SalaryMin, in.SalaryMax, in.Source,
		in.Priority, in.AppliedDate, in.Notes, id,
	)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

// UpdateStatus moves an application to a new status AND records the
// transition in application_status_history, inside one transaction. This is
// the piece that was deliberately left unwired when the migrations were
// written - nothing else is allowed to write current_status_id.
func (s *ApplicationStore) UpdateStatus(ctx context.Context, id uint64, statusID uint8, note *string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `UPDATE applications SET current_status_id = ? WHERE id = ?`, statusID, id)
	if err != nil {
		return err
	}
	if err := checkAffected(res); err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO application_status_history (application_id, status_id, notes) VALUES (?, ?, ?)`,
		id, statusID, note,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ApplicationStore) Delete(ctx context.Context, id uint64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM applications WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

type StatusHistoryEntry struct {
	ID        uint64  `json:"id"`
	StatusID  uint8   `json:"status_id"`
	ChangedAt string  `json:"changed_at"`
	Notes     *string `json:"notes,omitempty"`
}

func (s *ApplicationStore) History(ctx context.Context, applicationID uint64) ([]StatusHistoryEntry, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, status_id, changed_at, notes FROM application_status_history
		 WHERE application_id = ? ORDER BY changed_at ASC`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []StatusHistoryEntry
	for rows.Next() {
		var e StatusHistoryEntry
		var changedAt time.Time
		if err := rows.Scan(&e.ID, &e.StatusID, &changedAt, &e.Notes); err != nil {
			return nil, err
		}
		e.ChangedAt = changedAt.Format(time.RFC3339)
		out = append(out, e)
	}
	return out, rows.Err()
}

// ---- shared helpers used by every store in this package ----

// scannable is satisfied by both *sql.Row and *sql.Rows, so scan functions
// below work whether they're fed a single-row QueryRow result or one row
// out of a multi-row Query loop.
type scannable interface {
	Scan(dest ...any) error
}

func checkAffected(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func nullStringPtr(n sql.NullString) *string {
	if !n.Valid {
		return nil
	}
	return &n.String
}

func scanApplication(s scannable) (models.Application, error) {
	var a models.Application
	var jobDescription, jobPostingURL, location, workMode, source, notes sql.NullString
	var salaryMin, salaryMax sql.NullInt64
	var appliedDate sql.NullTime
	var archivedAt sql.NullTime

	err := s.Scan(
		&a.ID, &a.CompanyName, &a.PositionTitle, &jobDescription, &jobPostingURL, &location,
		&workMode, &salaryMin, &salaryMax, &source, &a.CurrentStatusID, &a.Priority, &appliedDate,
		&notes, &a.CreatedAt, &a.UpdatedAt, &archivedAt,
	)
	if err != nil {
		return a, err
	}

	a.JobDescription = nullStringPtr(jobDescription)
	a.JobPostingURL = nullStringPtr(jobPostingURL)
	a.Location = nullStringPtr(location)
	a.WorkMode = nullStringPtr(workMode)
	a.Source = nullStringPtr(source)
	a.Notes = nullStringPtr(notes)
	if salaryMin.Valid {
		v := uint32(salaryMin.Int64)
		a.SalaryMin = &v
	}
	if salaryMax.Valid {
		v := uint32(salaryMax.Int64)
		a.SalaryMax = &v
	}
	if appliedDate.Valid {
		v := appliedDate.Time.Format("2006-01-02")
		a.AppliedDate = &v
	}
	if archivedAt.Valid {
		a.ArchivedAt = &archivedAt.Time
	}

	return a, nil
}
