package server

import (
	"net/http"

	"debugo_back/routes"

	"github.com/gorilla/mux"
)

type Server struct {
	address string
}

func NewServer(address string) *Server {
	return &Server{address: address}
}

func (s *Server) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/hello", routes.HandleHello).
		Methods(http.MethodGet)

	return http.ListenAndServe(s.address, router)
}
