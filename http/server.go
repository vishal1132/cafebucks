package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/vishal1132/cafebucks/config"
)

// Server is the http server struct
type Server struct {
	Mux    *mux.Router
	Logger zerolog.Logger
}

// RunServer runs httpserver
func RunServer(cfg config.C, logger zerolog.Logger) (*Server, error) {

	// set up signal caching
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	logger.Info().
		Str("env", string(cfg.Env)).
		Str("log_level", cfg.LogLevel.String())

	_, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func() {
		sig := <-signalCh

		cancel()

		logger.Info().
			Str("signal", sig.String()).
			Msg("shutting down http server gracefully")
	}()

	s := Server{
		Mux:    mux.NewRouter(),
		Logger: logger,
	}

	httpSrvr := &http.Server{
		Handler:     s.Mux,
		ReadTimeout: 20 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	// creating a tcp listener
	socketAddr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	logger.Info().
		Str("addr", socketAddr).
		Msg("binding to TCP socket")

	// set up the network socket
	listener, err := net.Listen("tcp", socketAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to open HTTP socket: %w", err)
	}

	// signal handling / graceful shutdown goroutine
	serveStop, serverShutdown := make(chan struct{}), make(chan struct{})
	var serveErr, shutdownErr error

	// HTTP server parent goroutine
	go func() {
		defer close(serveStop)
		serveErr = httpSrvr.Serve(listener)
	}()

	// signal handling / graceful shutdown goroutine
	go func() {
		defer close(serverShutdown)
		sig := <-signalCh

		logger.Info().
			Str("signal", sig.String()).
			Msg("shutting HTTP server down gracefully")

		cctx, ccancel := context.WithTimeout(context.Background(), 25*time.Second)

		defer ccancel()
		defer cancel()

		if shutdownErr = httpSrvr.Shutdown(cctx); shutdownErr != nil {
			logger.Error().
				Err(shutdownErr).
				Msg("failed to gracefully shut down HTTP server")
		}
	}()

	// wait for it to die
	<-serverShutdown
	<-serveStop

	// log errors for informational purposes
	logger.Info().
		AnErr("serve_err", serveErr).
		AnErr("shutdown_err", shutdownErr).
		Msg("server shut down")

	return &s, nil
}