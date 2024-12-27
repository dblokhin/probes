package probes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type state struct {
	startup   atomic.Bool
	readiness atomic.Bool
	liveness  atomic.Bool
}

var (
	appState state

	// probe server things
	running atomic.Bool
	server  http.Server

	errInvalidPort    = errors.New("invalid http port")
	errAlreadyRunning = errors.New("http server is already running")
)

// RunServer initializes and runs the HTTP server for k8s probes on:
// `/startup`, `/ready` and `/live` for appropriate probes.
func RunServer(host string, port int) error {
	if err := validatePort(port); err != nil {
		return err
	}

	if !running.CompareAndSwap(false, true) {
		return errAlreadyRunning
	}

	defer func() { running.Store(false) }()

	address := host + ":" + strconv.Itoa(port)
	mux := http.NewServeMux()
	mux.HandleFunc("/startup", probeHandler(appState.isStarted))
	mux.HandleFunc("/ready", probeHandler(appState.isReady))
	mux.HandleFunc("/live", probeHandler(appState.isLive))

	const timeout = 3 * time.Second
	server = http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  -1,
	}

	appState.live(true)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown shuts down the HTTP server for probes.
func Shutdown(ctx context.Context) error {
	return server.Shutdown(ctx)
}

func (s *state) isStarted() bool {
	return s.startup.Load()
}

func (s *state) runStartup() {
	s.startup.CompareAndSwap(false, true)
}

func (s *state) isReady() bool {
	return s.readiness.Load()
}

func (s *state) ready(v bool) {
	s.startup.CompareAndSwap(false, true)
	s.readiness.Store(v)
}

func (s *state) isLive() bool {
	return s.liveness.Load()
}

func (s *state) live(v bool) {
	s.liveness.Store(v)
}

func validatePort(port int) error {
	if port > 0 && port < 65536 {
		return nil
	}

	return fmt.Errorf("%w: %q", errInvalidPort, port)
}

func probeHandler(stateFn func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if stateFn() {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}

// Startup updates the startup application state to true.
func Startup() {
	appState.runStartup()
}

// Ready marks the application as ready.
func Ready() {
	appState.ready(true)
}

// Unready marks the application as not ready.
func Unready() {
	appState.ready(false)
}

// Live marks the application as live.
func Live() {
	appState.live(true)
}

// Unlive marks the application as not live.
func Unlive() {
	appState.live(false)
}
