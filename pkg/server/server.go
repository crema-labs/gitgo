package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crema-labs/gitgo/pkg/store"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})

	// Create a new router
	s.router = mux.NewRouter()

	// Use the CORS middleware
	s.router = c.Handler(s.router).(*mux.Router)
	s.router.HandleFunc("/grant/{id}", s.handleGetGrant).Methods("GET")
	s.router.HandleFunc("/grant/{id}", s.handleUpdateGrantStatus).Methods("POST")
	s.router.HandleFunc("/grants", s.handleGetAllGrants).Methods("GET")
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

func (s *Server) handleGetAllGrants(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	grants, err := s.store.GetAllGrants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grants)
}
