package server

import (
	"diggle/searcher/handlers"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	handler *handlers.Handler
}

func NewServer(handler *handlers.Handler) *Server {
	return &Server{handler}
}

func (s *Server) Start() {
	router := http.NewServeMux()

	router.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	router.HandleFunc("GET /api/search", s.handler.HandleSearch)

	router.Handle("/", http.FileServer(http.Dir("./www")))

	fmt.Println("Server is listening on port 8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
