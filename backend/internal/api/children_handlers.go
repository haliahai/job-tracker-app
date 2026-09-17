package api

import (
	"net/http"
	"time"

	"jobtracker/internal/store"
)

// ---------- Contacts ----------

func (s *Server) handleListContacts(w http.ResponseWriter, r *http.Request) {
	appID, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	contacts, err := s.contacts.ListByApplication(r.Context(), appID)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, contacts)
}

type contactRequest struct {
	Name        string  `json:"name"`
	Role        string  `json:"role"`
	Company     *string `json:"company"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	LinkedInURL *string `json:"linkedin_url"`
	Notes       *string `json:"notes"`
}

func (s *Server) handleCreateContact(w http.ResponseWriter, r *http.Request) {
	appID, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req contactRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Role == "" {
		writeError(w, http.StatusBadRequest, "name and role are required")
		return
	}

	id, err := s.contacts.Create(r.Context(), store.CreateContactInput{
		ApplicationID: appID,
		Name:          req.Name,
		Role:          req.Role,
		Company:       req.Company,
		Email:         req.Email,
		Phone:         req.Phone,
		LinkedInURL:   req.LinkedInURL,
		Notes:         req.Notes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]uint64{"id": id})
}

func (s *Server) handleUpdateContact(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req contactRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Role == "" {
		writeError(w, http.StatusBadRequest, "name and role are required")
		return
	}

	err = s.contacts.Update(r.Context(), id, store.CreateContactInput{
		Name:        req.Name,
		Role:        req.Role,
		Company:     req.Company,
		Email:       req.Email,
		Phone:       req.Phone,
		LinkedInURL: req.LinkedInURL,
		Notes:       req.Notes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDeleteContact(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.contacts.Delete(r.Context(), id); err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---------- Interview stages ----------

func (s *Server) handleListInterviewStages(w http.ResponseWriter, r *http.Request) {
	appID, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	stages, err := s.interviewStages.ListByApplication(r.Context(), appID)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, stages)
}

type interviewStageRequest struct {
	StageName             string     `json:"stage_name"`
	StageType             string     `json:"stage_type"`
	ScheduledAt           *time.Time `json:"scheduled_at"`
	CompletedAt           *time.Time `json:"completed_at"`
	Format                *string    `json:"format"`
	InterviewerContactID  *uint64    `json:"interviewer_contact_id"`
	Outcome               string     `json:"outcome"`
	FeedbackNotes         *string    `json:"feedback_notes"`
}

func (s *Server) handleCreateInterviewStage(w http.ResponseWriter, r *http.Request) {
	appID, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req interviewStageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.StageName == "" || req.StageType == "" {
		writeError(w, http.StatusBadRequest, "stage_name and stage_type are required")
		return
	}
	if req.Outcome == "" {
		req.Outcome = "pending"
	}

	id, err := s.interviewStages.Create(r.Context(), store.UpsertInterviewStageInput{
		ApplicationID:        appID,
		StageName:            req.StageName,
		StageType:            req.StageType,
		ScheduledAt:          req.ScheduledAt,
		CompletedAt:          req.CompletedAt,
		Format:               req.Format,
		InterviewerContactID: req.InterviewerContactID,
		Outcome:              req.Outcome,
		FeedbackNotes:        req.FeedbackNotes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]uint64{"id": id})
}

func (s *Server) handleUpdateInterviewStage(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req interviewStageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.StageName == "" || req.StageType == "" {
		writeError(w, http.StatusBadRequest, "stage_name and stage_type are required")
		return
	}

	err = s.interviewStages.Update(r.Context(), id, store.UpsertInterviewStageInput{
		StageName:             req.StageName,
		StageType:             req.StageType,
		ScheduledAt:           req.ScheduledAt,
		CompletedAt:           req.CompletedAt,
		Format:                req.Format,
		InterviewerContactID:  req.InterviewerContactID,
		Outcome:               req.Outcome,
		FeedbackNotes:         req.FeedbackNotes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDeleteInterviewStage(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.interviewStages.Delete(r.Context(), id); err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---------- Prep items ----------

func (s *Server) handleListPrepItems(w http.ResponseWriter, r *http.Request) {
	appID, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	items, err := s.prepItems.ListByApplication(r.Context(), appID)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type prepItemRequest struct {
	InterviewStageID *uint64 `json:"interview_stage_id"`
	Category         string  `json:"category"`
	Kind             string  `json:"kind"`
	Content          string  `json:"content"`
	MyAnswerNotes    *string `json:"my_answer_notes"`
}

func (s *Server) handleCreatePrepItem(w http.ResponseWriter, r *http.Request) {
	appID, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req prepItemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}
	if req.Category == "" {
		req.Category = "other"
	}
	if req.Kind == "" {
		req.Kind = "to_prepare"
	}

	id, err := s.prepItems.Create(r.Context(), store.UpsertPrepItemInput{
		ApplicationID:    appID,
		InterviewStageID: req.InterviewStageID,
		Category:         req.Category,
		Kind:             req.Kind,
		Content:          req.Content,
		MyAnswerNotes:    req.MyAnswerNotes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]uint64{"id": id})
}

func (s *Server) handleUpdatePrepItem(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req prepItemRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	err = s.prepItems.Update(r.Context(), id, store.UpsertPrepItemInput{
		InterviewStageID: req.InterviewStageID,
		Category:         req.Category,
		Kind:             req.Kind,
		Content:          req.Content,
		MyAnswerNotes:    req.MyAnswerNotes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDeletePrepItem(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.prepItems.Delete(r.Context(), id); err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---------- Documents ----------

func (s *Server) handleListDocuments(w http.ResponseWriter, r *http.Request) {
	appID, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	docs, err := s.documents.ListByApplication(r.Context(), appID)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, docs)
}

type documentRequest struct {
	DocumentType string  `json:"document_type"`
	VersionLabel string  `json:"version_label"`
	FilePath     *string `json:"file_path"`
	Notes        *string `json:"notes"`
}

func (s *Server) handleCreateDocument(w http.ResponseWriter, r *http.Request) {
	appID, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req documentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.DocumentType == "" || req.VersionLabel == "" {
		writeError(w, http.StatusBadRequest, "document_type and version_label are required")
		return
	}

	id, err := s.documents.Create(r.Context(), store.UpsertDocumentInput{
		ApplicationID: appID,
		DocumentType:  req.DocumentType,
		VersionLabel:  req.VersionLabel,
		FilePath:      req.FilePath,
		Notes:         req.Notes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]uint64{"id": id})
}

func (s *Server) handleUpdateDocument(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req documentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.DocumentType == "" || req.VersionLabel == "" {
		writeError(w, http.StatusBadRequest, "document_type and version_label are required")
		return
	}

	err = s.documents.Update(r.Context(), id, store.UpsertDocumentInput{
		DocumentType: req.DocumentType,
		VersionLabel: req.VersionLabel,
		FilePath:     req.FilePath,
		Notes:        req.Notes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.documents.Delete(r.Context(), id); err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---------- Deadlines ----------

func (s *Server) handleListDeadlines(w http.ResponseWriter, r *http.Request) {
	appID, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	deadlines, err := s.deadlines.ListByApplication(r.Context(), appID)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, deadlines)
}

type deadlineRequest struct {
	DeadlineType string    `json:"deadline_type"`
	DueDate      time.Time `json:"due_date"`
	IsCompleted  bool      `json:"is_completed"`
	Notes        *string   `json:"notes"`
}

func (s *Server) handleCreateDeadline(w http.ResponseWriter, r *http.Request) {
	appID, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req deadlineRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.DeadlineType == "" {
		writeError(w, http.StatusBadRequest, "deadline_type is required")
		return
	}
	if req.DueDate.IsZero() {
		writeError(w, http.StatusBadRequest, "due_date is required")
		return
	}

	id, err := s.deadlines.Create(r.Context(), store.UpsertDeadlineInput{
		ApplicationID: appID,
		DeadlineType:  req.DeadlineType,
		DueDate:       req.DueDate,
		IsCompleted:   req.IsCompleted,
		Notes:         req.Notes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]uint64{"id": id})
}

func (s *Server) handleUpdateDeadline(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req deadlineRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.DeadlineType == "" {
		writeError(w, http.StatusBadRequest, "deadline_type is required")
		return
	}

	err = s.deadlines.Update(r.Context(), id, store.UpsertDeadlineInput{
		DeadlineType: req.DeadlineType,
		DueDate:      req.DueDate,
		IsCompleted:  req.IsCompleted,
		Notes:        req.Notes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDeleteDeadline(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.deadlines.Delete(r.Context(), id); err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
