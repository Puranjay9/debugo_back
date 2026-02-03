package routes

import (
	"debugo_back/db"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type HistoryObject struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Error       string    `json:"error"`
	ErrorSource string    `json:"error_source"`
	CodeDiff    string    `json:"code_diff"`
	Timestamp   time.Time `json:"timestamp"`
}

type RecordHistoryRequest struct {
	ProjectID   string `json:"project_id"`
	Error       string `json:"error"`
	ErrorSource string `json:"error_source"`
	CodeDiff    string `json:"code_diff"`
}

func HandleRecordApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var req RecordHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if req.ProjectID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "project_id is required",
		})
		return
	}

	history, err := createHistoryRecord(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create history record",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(history)
}

func createHistoryRecord(req RecordHistoryRequest) (*HistoryObject, error) {
	id := uuid.New().String()
	timestamp := time.Now()

	query := `INSERT INTO history (id, project_id, error, error_source, code_diff, timestamp) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.DB.Exec(query, id, req.ProjectID, req.Error, req.ErrorSource, req.CodeDiff, timestamp)
	if err != nil {
		return nil, err
	}

	return &HistoryObject{
		ID:          id,
		ProjectID:   req.ProjectID,
		Error:       req.Error,
		ErrorSource: req.ErrorSource,
		CodeDiff:    req.CodeDiff,
		Timestamp:   timestamp,
	}, nil
}
