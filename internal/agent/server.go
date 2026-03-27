package agent

import "net/http"

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/containers/run", s.handleRunContainer)

	return mux
}
