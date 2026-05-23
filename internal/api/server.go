package api

import (
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/frankfoo/doraemon/internal/store"
)

type Server struct {
	store     *store.Store
	jwtSecret []byte
	mux       *http.ServeMux
	addr      string
}

func New(s *store.Store, addr string, jwtSecret []byte, webFS fs.FS) *Server {
	srv := &Server{
		store:     s,
		jwtSecret: jwtSecret,
		addr:      addr,
		mux:       http.NewServeMux(),
	}
	srv.registerRoutes(webFS)
	return srv
}

func (s *Server) registerRoutes(webFS fs.FS) {
	s.mux.HandleFunc("GET /api/readme", s.handleReadme)
	s.mux.HandleFunc("GET /api/status", s.handleStatus)
	s.mux.HandleFunc("POST /api/setup", s.handleSetup)
	s.mux.HandleFunc("POST /api/login", s.handleLogin)

	s.mux.HandleFunc("GET /api/dashboard", s.authRequired(s.handleDashboard))

	s.mux.HandleFunc("GET /api/entries/search", s.authRequired(s.handleSearchEntries))
	s.mux.HandleFunc("GET /api/entries/{id}", s.authRequired(s.handleGetEntry))
	s.mux.HandleFunc("GET /api/entries", s.authRequired(s.handleListEntries))
	s.mux.HandleFunc("POST /api/entries", s.authRequired(s.handleCreateEntry))
	s.mux.HandleFunc("PUT /api/entries/{id}", s.authRequired(s.handleUpdateEntry))
	s.mux.HandleFunc("DELETE /api/entries/{id}", s.authRequired(s.handleDeleteEntry))

	if webFS != nil {
		s.mux.Handle("/", spaHandler(webFS))
	}
}

func (s *Server) ListenAndServe() error {
	handler := corsMiddleware(loggingMiddleware(s.mux))
	return http.ListenAndServe(s.addr, handler)
}

func (s *Server) Addr() string {
	return s.addr
}

func spaHandler(webFS fs.FS) http.Handler {
	fileServer := http.FileServerFS(webFS)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else {
			path = path[1:] // strip leading slash
		}
		if _, err := fs.Stat(webFS, path); err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
