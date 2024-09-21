package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crema-labs/gitgo/pkg/store"
	"github.com/gorilla/mux"
)

type Server struct {
	store  store.Store
	router *mux.Router
}

func NewServer(store store.Store) *Server {
	s := &Server{
		store:  store,
		router: mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/grant/{id}", s.handleGetGrant).Methods("GET")
	s.router.HandleFunc("/grant/{id}", s.handleUpdateGrantStatus).Methods("POST")
}

func (s *Server) Start(addr string) error {
	fmt.Printf("Server is running on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) handleGetGrant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	grantID := vars["id"]

	grant, err := s.store.GetGrant(grantID)
	if err != nil {
		if err == store.ErrGrantNotFound {
			http.Error(w, "Grant not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grant)
}

func (s *Server) handleUpdateGrantStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	grantID := vars["id"]

	err := s.store.UpdateGrantStatus(grantID, "closed")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Grant status updated to closed")
}
