package httpserv

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/freundallein/resender/producers"
	"github.com/freundallein/resender/uidgen"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// ErrNoOptions - you should provide Options
	ErrNoOptions = errors.New("no options provided")
)

// Options - server's parameters
type Options struct {
	Port      string
	Gen       *uidgen.Generator
	Producers []producers.Producer
}

// Server - main server struct
type Server struct {
	options *Options
}

// New - server constructor
func New(options *Options) (*Server, error) {
	if options == nil {
		return nil, ErrNoOptions
	}
	return &Server{options: options}, nil
}

// Run - start server
func (srv *Server) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", Index(srv.options))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	addr := fmt.Sprintf("0.0.0.0:%s", srv.options.Port)
	serv := &http.Server{
		Handler:        mux,
		Addr:           addr,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("[server] start listening on :%s\n", srv.options.Port)
	err := serv.ListenAndServe()
	return err
}
