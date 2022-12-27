package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/obliviousfrog/flighttracker/internal/tracker"

	"go.uber.org/zap"
)

// Config contains the configurations needed to run the server
type Config struct {
	Host    string
	Port    int
	Log     *zap.Logger
	Tracker *tracker.Tracker
}

// Server -
type Server struct {
	host    string
	port    int
	tracker *tracker.Tracker
	log     *zap.Logger
	router  *mux.Router
}

// New returns a newly configured server
func New(cfg Config) *Server {
	server := Server{
		host:    cfg.Host,
		port:    cfg.Port,
		log:     cfg.Log,
		tracker: cfg.Tracker,
	}

	r := mux.NewRouter()
	r.HandleFunc("/calculate", server.calculateFlights).Methods("POST")
	r.Use(server.loggingMiddleware)
	server.router = r

	return &server
}

// Start starts the http server
func (s *Server) Start() error {
	srv := &http.Server{
		Handler:      s.router,
		Addr:         fmt.Sprintf("%s:%d", s.host, s.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) loggingMiddleware(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(rw, r)

		duration := time.Since(start)
		s.log.Info(
			"Processed Request.",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Duration("duration", duration),
		)
	}
	return http.HandlerFunc(logFn)
}

func (s *Server) calculateFlights(w http.ResponseWriter, r *http.Request) {
	var flights [][]string
	err := json.NewDecoder(r.Body).Decode(&flights)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	flight, err := s.tracker.GetSrcAndDst(tracker.Flights(flights))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&flight); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
