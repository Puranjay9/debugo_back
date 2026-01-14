package server

import (
	"debugo_back/db"
	"debugo_back/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	address string
}

func NewServer(address string) *Server {
	return &Server{address: address}
}

func (s *Server) Run() error {

	//dB connection
	if err := db.ConnectToDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.CloseDB()

	//router configs
	router := mux.NewRouter()

	router.HandleFunc("/cli/init", routes.HandleProjectInit).
		Methods(http.MethodPost)

	return http.ListenAndServe(s.address, router)
}
