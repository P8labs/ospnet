package master

import "net/http"

type Server struct {
	store *Store
}

func NewServer() *Server {
	return &Server{
		store: NewStore(),
	}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/nodes/register", s.handleRegister)
	mux.HandleFunc("/nodes/heartbeat", s.handleHeartbeat)
	mux.HandleFunc("/nodes", s.handleListNodes)

	return mux
}
