package nodeagent

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime"
	"time"
)

type Server struct {
	cfg    *Config
	logger *slog.Logger
	mux    *http.ServeMux
}

func NewServer(cfg *Config, logger *slog.Logger) (*Server, error) {
	s := &Server{
		cfg:    cfg,
		logger: logger,
		mux:    http.NewServeMux(),
	}

	s.routes()
	return s, nil
}

func (s *Server) Router() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.healthz)
	s.mux.HandleFunc("/inventory", s.inventory)
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"agentId": s.cfg.AgentID,
		"hostId": s.cfg.HostID,
		"time": time.Now().UTC(),
	})
}

func (s *Server) inventory(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]any{
		"agentId": s.cfg.AgentID,
		"hostId": s.cfg.HostID,
		"goVersion": runtime.Version(),
		"os": runtime.GOOS,
		"arch": runtime.GOARCH,
		"time": time.Now().UTC(),
	})
}
