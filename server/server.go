package server

import (
	"net/http"

	"debugo_back/routes"

	"github.com/gorilla/mux"

	"debugo_back/db"

	"log"
)

type Server struct {
	address string
}

func NewServer(address string) *Server {
	return &Server{address: address}
}

func (s *Server) Run() error {

	if err := db.ConnectToDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.CloseDB()
	router := mux.NewRouter()

	router.HandleFunc("/hello", routes.HandleHello).
		Methods(http.MethodGet)

	return http.ListenAndServe(s.address, router)
}
