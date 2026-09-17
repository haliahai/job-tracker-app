package api

import (
	"database/sql"
	"net/http"

	"jobtracker/internal/store"
)

type Server struct {
	statuses        *store.StatusStore
	applications    *store.ApplicationStore
	contacts        *store.ContactStore
	interviewStages *store.InterviewStageStore
	prepItems       *store.PrepItemStore
	documents       *store.DocumentStore
	deadlines       *store.DeadlineStore
}

func NewRouter(db *sql.DB) http.Handler {
	s := &Server{
		statuses:        store.NewStatusStore(db),
		applications:    store.NewApplicationStore(db),
		contacts:        store.NewContactStore(db),
		interviewStages: store.NewInterviewStageStore(db),
		prepItems:       store.NewPrepItemStore(db),
		documents:       store.NewDocumentStore(db),
		deadlines:       store.NewDeadlineStore(db),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", s.handleHealth)

	mux.HandleFunc("GET /api/statuses", s.handleListStatuses)

	mux.HandleFunc("GET /api/applications", s.handleListApplications)
	mux.HandleFunc("POST /api/applications", s.handleCreateApplication)
	mux.HandleFunc("GET /api/applications/{id}", s.handleGetApplication)
	mux.HandleFunc("PUT /api/applications/{id}", s.handleUpdateApplication)
	mux.HandleFunc("PATCH /api/applications/{id}/status", s.handleUpdateApplicationStatus)
	mux.HandleFunc("DELETE /api/applications/{id}", s.handleDeleteApplication)
	mux.HandleFunc("GET /api/applications/{id}/history", s.handleApplicationHistory)

	mux.HandleFunc("GET /api/applications/{id}/contacts", s.handleListContacts)
	mux.HandleFunc("POST /api/applications/{id}/contacts", s.handleCreateContact)
	mux.HandleFunc("PUT /api/contacts/{id}", s.handleUpdateContact)
	mux.HandleFunc("DELETE /api/contacts/{id}", s.handleDeleteContact)

	mux.HandleFunc("GET /api/applications/{id}/interview-stages", s.handleListInterviewStages)
	mux.HandleFunc("POST /api/applications/{id}/interview-stages", s.handleCreateInterviewStage)
	mux.HandleFunc("PUT /api/interview-stages/{id}", s.handleUpdateInterviewStage)
	mux.HandleFunc("DELETE /api/interview-stages/{id}", s.handleDeleteInterviewStage)

	mux.HandleFunc("GET /api/applications/{id}/prep-items", s.handleListPrepItems)
	mux.HandleFunc("POST /api/applications/{id}/prep-items", s.handleCreatePrepItem)
	mux.HandleFunc("PUT /api/prep-items/{id}", s.handleUpdatePrepItem)
	mux.HandleFunc("DELETE /api/prep-items/{id}", s.handleDeletePrepItem)

	mux.HandleFunc("GET /api/applications/{id}/documents", s.handleListDocuments)
	mux.HandleFunc("POST /api/applications/{id}/documents", s.handleCreateDocument)
	mux.HandleFunc("PUT /api/documents/{id}", s.handleUpdateDocument)
	mux.HandleFunc("DELETE /api/documents/{id}", s.handleDeleteDocument)

	mux.HandleFunc("GET /api/applications/{id}/deadlines", s.handleListDeadlines)
	mux.HandleFunc("POST /api/applications/{id}/deadlines", s.handleCreateDeadline)
	mux.HandleFunc("PUT /api/deadlines/{id}", s.handleUpdateDeadline)
	mux.HandleFunc("DELETE /api/deadlines/{id}", s.handleDeleteDeadline)

	return recoverMiddleware(loggingMiddleware(corsMiddleware(mux)))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
