package store

import (
	"context"
	"database/sql"
	"time"

	"jobtracker/internal/models"
)

// ---------- Contacts ----------

type ContactStore struct{ db *sql.DB }

func NewContactStore(db *sql.DB) *ContactStore { return &ContactStore{db: db} }

func (s *ContactStore) ListByApplication(ctx context.Context, applicationID uint64) ([]models.Contact, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, application_id, name, role, company, email, phone, linkedin_url, notes, created_at
		 FROM contacts WHERE application_id = ? ORDER BY created_at DESC`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Contact
	for rows.Next() {
		c, err := scanContact(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

type CreateContactInput struct {
	ApplicationID uint64
	Name          string
	Role          string
	Company       *string
	Email         *string
	Phone         *string
	LinkedInURL   *string
	Notes         *string
}

func (s *ContactStore) Create(ctx context.Context, in CreateContactInput) (uint64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO contacts (application_id, name, role, company, email, phone, linkedin_url, notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		in.ApplicationID, in.Name, in.Role, in.Company, in.Email, in.Phone, in.LinkedInURL, in.Notes,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (s *ContactStore) Update(ctx context.Context, id uint64, in CreateContactInput) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE contacts SET name = ?, role = ?, company = ?, email = ?, phone = ?, linkedin_url = ?, notes = ?
		 WHERE id = ?`,
		in.Name, in.Role, in.Company, in.Email, in.Phone, in.LinkedInURL, in.Notes, id,
	)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

func (s *ContactStore) Delete(ctx context.Context, id uint64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM contacts WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

func scanContact(s scannable) (models.Contact, error) {
	var c models.Contact
	var company, email, phone, linkedin, notes sql.NullString
	err := s.Scan(&c.ID, &c.ApplicationID, &c.Name, &c.Role, &company, &email, &phone, &linkedin, &notes, &c.CreatedAt)
	if err != nil {
		return c, err
	}
	c.Company = nullStringPtr(company)
	c.Email = nullStringPtr(email)
	c.Phone = nullStringPtr(phone)
	c.LinkedInURL = nullStringPtr(linkedin)
	c.Notes = nullStringPtr(notes)
	return c, nil
}

// ---------- Interview stages ----------

type InterviewStageStore struct{ db *sql.DB }

func NewInterviewStageStore(db *sql.DB) *InterviewStageStore { return &InterviewStageStore{db: db} }

func (s *InterviewStageStore) ListByApplication(ctx context.Context, applicationID uint64) ([]models.InterviewStage, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, application_id, stage_name, stage_type, scheduled_at, completed_at, format,
		        interviewer_contact_id, outcome, feedback_notes, created_at, updated_at
		 FROM interview_stages WHERE application_id = ? ORDER BY COALESCE(scheduled_at, created_at) ASC`,
		applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.InterviewStage
	for rows.Next() {
		st, err := scanInterviewStage(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, rows.Err()
}

type UpsertInterviewStageInput struct {
	ApplicationID        uint64
	StageName            string
	StageType            string
	ScheduledAt          *time.Time
	CompletedAt          *time.Time
	Format               *string
	InterviewerContactID *uint64
	Outcome              string
	FeedbackNotes        *string
}

func (s *InterviewStageStore) Create(ctx context.Context, in UpsertInterviewStageInput) (uint64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO interview_stages
			(application_id, stage_name, stage_type, scheduled_at, completed_at, format,
			 interviewer_contact_id, outcome, feedback_notes)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		in.ApplicationID, in.StageName, in.StageType, in.ScheduledAt, in.CompletedAt, in.Format,
		in.InterviewerContactID, in.Outcome, in.FeedbackNotes,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (s *InterviewStageStore) Update(ctx context.Context, id uint64, in UpsertInterviewStageInput) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE interview_stages SET
			stage_name = ?, stage_type = ?, scheduled_at = ?, completed_at = ?, format = ?,
			interviewer_contact_id = ?, outcome = ?, feedback_notes = ?
		 WHERE id = ?`,
		in.StageName, in.StageType, in.ScheduledAt, in.CompletedAt, in.Format,
		in.InterviewerContactID, in.Outcome, in.FeedbackNotes, id,
	)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

func (s *InterviewStageStore) Delete(ctx context.Context, id uint64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM interview_stages WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

func scanInterviewStage(s scannable) (models.InterviewStage, error) {
	var st models.InterviewStage
	var scheduledAt, completedAt sql.NullTime
	var format, feedbackNotes sql.NullString
	var interviewerID sql.NullInt64

	err := s.Scan(&st.ID, &st.ApplicationID, &st.StageName, &st.StageType, &scheduledAt, &completedAt,
		&format, &interviewerID, &st.Outcome, &feedbackNotes, &st.CreatedAt, &st.UpdatedAt)
	if err != nil {
		return st, err
	}
	if scheduledAt.Valid {
		st.ScheduledAt = &scheduledAt.Time
	}
	if completedAt.Valid {
		st.CompletedAt = &completedAt.Time
	}
	st.Format = nullStringPtr(format)
	st.FeedbackNotes = nullStringPtr(feedbackNotes)
	if interviewerID.Valid {
		v := uint64(interviewerID.Int64)
		st.InterviewerContactID = &v
	}
	return st, nil
}

// ---------- Prep items ----------

type PrepItemStore struct{ db *sql.DB }

func NewPrepItemStore(db *sql.DB) *PrepItemStore { return &PrepItemStore{db: db} }

func (s *PrepItemStore) ListByApplication(ctx context.Context, applicationID uint64) ([]models.PrepItem, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, application_id, interview_stage_id, category, kind, content, my_answer_notes, created_at
		 FROM prep_items WHERE application_id = ? ORDER BY created_at DESC`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.PrepItem
	for rows.Next() {
		p, err := scanPrepItem(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

type UpsertPrepItemInput struct {
	ApplicationID    uint64
	InterviewStageID *uint64
	Category         string
	Kind             string
	Content          string
	MyAnswerNotes    *string
}

func (s *PrepItemStore) Create(ctx context.Context, in UpsertPrepItemInput) (uint64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO prep_items (application_id, interview_stage_id, category, kind, content, my_answer_notes)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		in.ApplicationID, in.InterviewStageID, in.Category, in.Kind, in.Content, in.MyAnswerNotes,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (s *PrepItemStore) Update(ctx context.Context, id uint64, in UpsertPrepItemInput) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE prep_items SET interview_stage_id = ?, category = ?, kind = ?, content = ?, my_answer_notes = ?
		 WHERE id = ?`,
		in.InterviewStageID, in.Category, in.Kind, in.Content, in.MyAnswerNotes, id,
	)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

func (s *PrepItemStore) Delete(ctx context.Context, id uint64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM prep_items WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

func scanPrepItem(s scannable) (models.PrepItem, error) {
	var p models.PrepItem
	var stageID sql.NullInt64
	var answerNotes sql.NullString

	err := s.Scan(&p.ID, &p.ApplicationID, &stageID, &p.Category, &p.Kind, &p.Content, &answerNotes, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	if stageID.Valid {
		v := uint64(stageID.Int64)
		p.InterviewStageID = &v
	}
	p.MyAnswerNotes = nullStringPtr(answerNotes)
	return p, nil
}

// ---------- Documents ----------

type DocumentStore struct{ db *sql.DB }

func NewDocumentStore(db *sql.DB) *DocumentStore { return &DocumentStore{db: db} }

func (s *DocumentStore) ListByApplication(ctx context.Context, applicationID uint64) ([]models.Document, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, application_id, document_type, version_label, file_path, notes, created_at
		 FROM documents WHERE application_id = ? ORDER BY created_at DESC`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Document
	for rows.Next() {
		d, err := scanDocument(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

type UpsertDocumentInput struct {
	ApplicationID uint64
	DocumentType  string
	VersionLabel  string
	FilePath      *string
	Notes         *string
}

func (s *DocumentStore) Create(ctx context.Context, in UpsertDocumentInput) (uint64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO documents (application_id, document_type, version_label, file_path, notes)
		 VALUES (?, ?, ?, ?, ?)`,
		in.ApplicationID, in.DocumentType, in.VersionLabel, in.FilePath, in.Notes,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (s *DocumentStore) Update(ctx context.Context, id uint64, in UpsertDocumentInput) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE documents SET document_type = ?, version_label = ?, file_path = ?, notes = ? WHERE id = ?`,
		in.DocumentType, in.VersionLabel, in.FilePath, in.Notes, id,
	)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

func (s *DocumentStore) Delete(ctx context.Context, id uint64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM documents WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

func scanDocument(s scannable) (models.Document, error) {
	var d models.Document
	var filePath, notes sql.NullString
	err := s.Scan(&d.ID, &d.ApplicationID, &d.DocumentType, &d.VersionLabel, &filePath, &notes, &d.CreatedAt)
	if err != nil {
		return d, err
	}
	d.FilePath = nullStringPtr(filePath)
	d.Notes = nullStringPtr(notes)
	return d, nil
}

// ---------- Deadlines ----------

type DeadlineStore struct{ db *sql.DB }

func NewDeadlineStore(db *sql.DB) *DeadlineStore { return &DeadlineStore{db: db} }

func (s *DeadlineStore) ListByApplication(ctx context.Context, applicationID uint64) ([]models.Deadline, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, application_id, deadline_type, due_date, is_completed, notes, created_at
		 FROM deadlines WHERE application_id = ? ORDER BY due_date ASC`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Deadline
	for rows.Next() {
		dl, err := scanDeadline(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, dl)
	}
	return out, rows.Err()
}

type UpsertDeadlineInput struct {
	ApplicationID uint64
	DeadlineType  string
	DueDate       time.Time
	IsCompleted   bool
	Notes         *string
}

func (s *DeadlineStore) Create(ctx context.Context, in UpsertDeadlineInput) (uint64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO deadlines (application_id, deadline_type, due_date, is_completed, notes)
		 VALUES (?, ?, ?, ?, ?)`,
		in.ApplicationID, in.DeadlineType, in.DueDate, in.IsCompleted, in.Notes,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (s *DeadlineStore) Update(ctx context.Context, id uint64, in UpsertDeadlineInput) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE deadlines SET deadline_type = ?, due_date = ?, is_completed = ?, notes = ? WHERE id = ?`,
		in.DeadlineType, in.DueDate, in.IsCompleted, in.Notes, id,
	)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

func (s *DeadlineStore) Delete(ctx context.Context, id uint64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM deadlines WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return checkAffected(res)
}

func scanDeadline(s scannable) (models.Deadline, error) {
	var dl models.Deadline
	var notes sql.NullString
	err := s.Scan(&dl.ID, &dl.ApplicationID, &dl.DeadlineType, &dl.DueDate, &dl.IsCompleted, &notes, &dl.CreatedAt)
	if err != nil {
		return dl, err
	}
	dl.Notes = nullStringPtr(notes)
	return dl, nil
}
