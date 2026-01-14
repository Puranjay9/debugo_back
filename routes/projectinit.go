package routes

import (
	"debugo_back/db"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type InitProjectRequest struct {
	Name string `json:"name"`
}

type InitProjectResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func HandleProjectInit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req InitProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Request Body",
		})
		return
	}

	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Project Name is Required",
		})
		return
	}

	project, err := createProject(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create project",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(InitProjectResponse{
		ID:   project.ID,
		Name: project.Name,
	})

}

func createProject(name string) (*Project, error) {
	id := uuid.New().String()

	query := `INSERT INTO projects (id, name) VALUES ($1, $2)`
	_, err := db.DB.Exec(query, id, name)
	if err != nil {
		return nil, err
	}

	return &Project{
		ID:   id,
		Name: name,
	}, nil
}
