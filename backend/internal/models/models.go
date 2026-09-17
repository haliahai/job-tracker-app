package models

import "time"

type ApplicationStatus struct {
	ID        uint8  `json:"id"`
	Code      string `json:"code"`
	Label     string `json:"label"`
	SortOrder uint8  `json:"sort_order"`
}

type Application struct {
	ID              uint64     `json:"id"`
	CompanyName     string     `json:"company_name"`
	PositionTitle   string     `json:"position_title"`
	JobDescription  *string    `json:"job_description,omitempty"`
	JobPostingURL   *string    `json:"job_posting_url,omitempty"`
	Location        *string    `json:"location,omitempty"`
	WorkMode        *string    `json:"work_mode,omitempty"`
	SalaryMin       *uint32    `json:"salary_min,omitempty"`
	SalaryMax       *uint32    `json:"salary_max,omitempty"`
	Source          *string    `json:"source,omitempty"`
	CurrentStatusID uint8      `json:"current_status_id"`
	Priority        string     `json:"priority"`
	AppliedDate     *string    `json:"applied_date,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ArchivedAt      *time.Time `json:"archived_at,omitempty"`
}

type Contact struct {
	ID            uint64    `json:"id"`
	ApplicationID uint64    `json:"application_id"`
	Name          string    `json:"name"`
	Role          string    `json:"role"`
	Company       *string   `json:"company,omitempty"`
	Email         *string   `json:"email,omitempty"`
	Phone         *string   `json:"phone,omitempty"`
	LinkedInURL   *string   `json:"linkedin_url,omitempty"`
	Notes         *string   `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type InterviewStage struct {
	ID                    uint64     `json:"id"`
	ApplicationID         uint64     `json:"application_id"`
	StageName             string     `json:"stage_name"`
	StageType             string     `json:"stage_type"`
	ScheduledAt           *time.Time `json:"scheduled_at,omitempty"`
	CompletedAt           *time.Time `json:"completed_at,omitempty"`
	Format                *string    `json:"format,omitempty"`
	InterviewerContactID  *uint64    `json:"interviewer_contact_id,omitempty"`
	Outcome               string     `json:"outcome"`
	FeedbackNotes         *string    `json:"feedback_notes,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type PrepItem struct {
	ID               uint64    `json:"id"`
	ApplicationID    uint64    `json:"application_id"`
	InterviewStageID *uint64   `json:"interview_stage_id,omitempty"`
	Category         string    `json:"category"`
	Kind             string    `json:"kind"`
	Content          string    `json:"content"`
	MyAnswerNotes    *string   `json:"my_answer_notes,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type Document struct {
	ID            uint64    `json:"id"`
	ApplicationID uint64    `json:"application_id"`
	DocumentType  string    `json:"document_type"`
	VersionLabel  string    `json:"version_label"`
	FilePath      *string   `json:"file_path,omitempty"`
	Notes         *string   `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type Deadline struct {
	ID            uint64    `json:"id"`
	ApplicationID uint64    `json:"application_id"`
	DeadlineType  string    `json:"deadline_type"`
	DueDate       time.Time `json:"due_date"`
	IsCompleted   bool      `json:"is_completed"`
	Notes         *string   `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
