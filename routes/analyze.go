package routes

import (
	"bytes"
	"database/sql"
	"debugo_back/db"
	"encoding/json"
	"fmt"
	"net/http"
)

// RagEndpoint is the URL of the RAG service
// TODO: Replace with actual RAG endpoint
const RagEndpoint = "http://localhost:8001/rag/analyze"

type AnalyzeRequest struct {
	HistoryID string `json:"history_id"`
}

type RagPayload struct {
	Error       string `json:"error"`
	ErrorSource string `json:"error_source"`
	CodeDiff    string `json:"code_diff"`
}

type RagResponse struct {
	Analysis string     `json:"analysis"`
	Fix      string     `json:"fix"`
	Payload  RagPayload `json:"payload,omitempty"`
}

func HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if req.HistoryID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "history_id is required",
		})
		return
	}

	// Fetch history record
	history, err := getHistoryByID(req.HistoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "History record not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Database error: " + err.Error(),
		})
		return
	}

	// Prepare payload for RAG
	ragPayload := RagPayload{
		Error:       history.Error,
		ErrorSource: history.ErrorSource,
		CodeDiff:    history.CodeDiff,
	}

	payloadBytes, err := json.Marshal(ragPayload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to serialize RAG payload",
		})
		return
	}

	// Hit RAG endpoint
	// Using a dummy request structure for now
	resp, err := http.Post(RagEndpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		// Mock response for testing when RAG is unavailable
		w.WriteHeader(http.StatusOK) // Returning 200 OK so you can see the payload
		json.NewEncoder(w).Encode(RagResponse{
			Analysis: "RAG Service Unavailable (Mock Response)",
			Fix:      "Check if RAG service is running",
			Payload:  ragPayload,
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("RAG service returned status %d", resp.StatusCode),
		})
		return
	}

	var ragResp RagResponse
	if err := json.NewDecoder(resp.Body).Decode(&ragResp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to decode RAG response",
		})
		return
	}

	ragResp.Payload = ragPayload // Attach payload to success response too

	// Return analysis to user
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ragResp)
}

func getHistoryByID(id string) (*HistoryObject, error) {
	var h HistoryObject
	query := `SELECT id, project_id, error, error_source, code_diff, timestamp FROM history WHERE id = $1`
	err := db.DB.QueryRow(query, id).Scan(&h.ID, &h.ProjectID, &h.Error, &h.ErrorSource, &h.CodeDiff, &h.Timestamp)
	if err != nil {
		return nil, err
	}
	return &h, nil
}
