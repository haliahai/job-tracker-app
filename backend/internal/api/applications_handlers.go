package api

import (
	"net/http"

	"jobtracker/internal/store"
)

func (s *Server) handleListApplications(w http.ResponseWriter, r *http.Request) {
	var filter store.ApplicationFilter

	if statusStr := r.URL.Query().Get("status_id"); statusStr != "" {
		id, err := parseUint8(statusStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid status_id")
			return
		}
		filter.StatusID = &id
	}
	filter.Company = r.URL.Query().Get("company")

	apps, err := s.applications.List(r.Context(), filter)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apps)
}

type applicationRequest struct {
	CompanyName     string  `json:"company_name"`
	PositionTitle   string  `json:"position_title"`
	JobDescription  *string `json:"job_description"`
	JobPostingURL   *string `json:"job_posting_url"`
	Location        *string `json:"location"`
	WorkMode        *string `json:"work_mode"`
	SalaryMin       *uint32 `json:"salary_min"`
	SalaryMax       *uint32 `json:"salary_max"`
	Source          *string `json:"source"`
	CurrentStatusID uint8   `json:"current_status_id"`
	Priority        string  `json:"priority"`
	AppliedDate     *string `json:"applied_date"`
	Notes           *string `json:"notes"`
}

func (req applicationRequest) valid() string {
	if req.CompanyName == "" {
		return "company_name is required"
	}
	if req.PositionTitle == "" {
		return "position_title is required"
	}
	if req.Priority == "" {
		return "priority is required"
	}
	return ""
}

func (s *Server) handleCreateApplication(w http.ResponseWriter, r *http.Request) {
	var req applicationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CurrentStatusID == 0 {
		writeError(w, http.StatusBadRequest, "current_status_id is required")
		return
	}
	if msg := req.valid(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	id, err := s.applications.Create(r.Context(), store.CreateApplicationInput{
		CompanyName:     req.CompanyName,
		PositionTitle:   req.PositionTitle,
		JobDescription:  req.JobDescription,
		JobPostingURL:   req.JobPostingURL,
		Location:        req.Location,
		WorkMode:        req.WorkMode,
		SalaryMin:       req.SalaryMin,
		SalaryMax:       req.SalaryMax,
		Source:          req.Source,
		CurrentStatusID: req.CurrentStatusID,
		Priority:        req.Priority,
		AppliedDate:     req.AppliedDate,
		Notes:           req.Notes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}

	app, err := s.applications.Get(r.Context(), id)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, app)
}

func (s *Server) handleGetApplication(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	app, err := s.applications.Get(r.Context(), id)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, app)
}

func (s *Server) handleUpdateApplication(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req applicationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if msg := req.valid(); msg != "" {
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	err = s.applications.Update(r.Context(), id, store.UpdateApplicationInput{
		CompanyName:    req.CompanyName,
		PositionTitle:  req.PositionTitle,
		JobDescription: req.JobDescription,
		JobPostingURL:  req.JobPostingURL,
		Location:       req.Location,
		WorkMode:       req.WorkMode,
		SalaryMin:      req.SalaryMin,
		SalaryMax:      req.SalaryMax,
		Source:         req.Source,
		Priority:       req.Priority,
		AppliedDate:    req.AppliedDate,
		Notes:          req.Notes,
	})
	if err != nil {
		handleStoreError(w, err)
		return
	}

	app, err := s.applications.Get(r.Context(), id)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, app)
}

type updateStatusRequest struct {
	StatusID uint8   `json:"status_id"`
	Notes    *string `json:"notes"`
}

func (s *Server) handleUpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req updateStatusRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.StatusID == 0 {
		writeError(w, http.StatusBadRequest, "status_id is required")
		return
	}

	if err := s.applications.UpdateStatus(r.Context(), id, req.StatusID, req.Notes); err != nil {
		handleStoreError(w, err)
		return
	}

	app, err := s.applications.Get(r.Context(), id)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, app)
}

func (s *Server) handleDeleteApplication(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := s.applications.Delete(r.Context(), id); err != nil {
		handleStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleApplicationHistory(w http.ResponseWriter, r *http.Request) {
	id, err := idParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	history, err := s.applications.History(r.Context(), id)
	if err != nil {
		handleStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, history)
}
