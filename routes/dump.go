package routes

import (
	"debugo_back/db"
	"encoding/json"
	"net/http"
)

type DatabaseDump struct {
	Projects []Project       `json:"projects"`
	History  []HistoryObject `json:"history"`
}

func HandleGetAllData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	dump, err := getAllData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to fetch data: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dump)
}

func getAllData() (*DatabaseDump, error) {
	dump := &DatabaseDump{}

	// Fetch all projects
	projectRows, err := db.DB.Query(`SELECT id, name FROM projects`)
	if err != nil {
		return nil, err
	}
	defer projectRows.Close()

	for projectRows.Next() {
		var p Project
		if err := projectRows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		dump.Projects = append(dump.Projects, p)
	}

	// Fetch all history
	historyRows, err := db.DB.Query(`SELECT id, project_id, error, error_source, code_diff, timestamp FROM history`)
	if err != nil {
		return nil, err
	}
	defer historyRows.Close()

	for historyRows.Next() {
		var h HistoryObject
		if err := historyRows.Scan(&h.ID, &h.ProjectID, &h.Error, &h.ErrorSource, &h.CodeDiff, &h.Timestamp); err != nil {
			return nil, err
		}
		dump.History = append(dump.History, h)
	}

	return dump, nil
}
