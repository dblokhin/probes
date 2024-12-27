package probes

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestState_isStarted tests the isStarted method of the state struct.
func TestState_isStarted(t *testing.T) {
	var s state
	if s.isStarted() {
		t.Error("expected isStarted to return false initially")
	}

	s.startup.Store(true)
	if !s.isStarted() {
		t.Error("expected isStarted to return true after being set")
	}
}

// TestState_startup tests the startup method of the state struct.
func TestState_startup(t *testing.T) {
	var s state
	s.runStartup()
	if !s.startup.Load() {
		t.Error("expected startup to be true after startup call")
	}

	// Ensure startup doesn't overwrite an already true value.
	s.runStartup()
	if !s.startup.Load() {
		t.Error("expected startup to remain true after second startup call")
	}
}

// TestState_isReady tests the isReady method of the state struct.
func TestState_isReady(t *testing.T) {
	var s state
	if s.isReady() {
		t.Error("expected isReady to return false initially")
	}

	s.readiness.Store(true)
	if !s.isReady() {
		t.Error("expected isReady to return true after being set")
	}
}

// TestState_ready tests the ready method of the state struct.
func TestState_ready(t *testing.T) {
	var s state
	s.ready(true)
	if !s.readiness.Load() {
		t.Error("expected readiness to be true after ready(true) call")
	}

	s.ready(false)
	if s.readiness.Load() {
		t.Error("expected readiness to be false after ready(false) call")
	}
}

// TestState_isLive tests the isLive method of the state struct.
func TestState_isLive(t *testing.T) {
	var s state
	if s.isLive() {
		t.Error("expected isLive to return false initially")
	}

	s.liveness.Store(true)
	if !s.isLive() {
		t.Error("expected isLive to return true after being set")
	}
}

// TestState_live tests the live method of the state struct.
func TestState_live(t *testing.T) {
	var s state
	s.live(true)
	if !s.liveness.Load() {
		t.Error("expected liveness to be true after live(true) call")
	}

	s.live(false)
	if s.liveness.Load() {
		t.Error("expected liveness to be false after live(false) call")
	}
}

// TestValidatePort tests the validatePort function.
func TestValidatePort(t *testing.T) {
	tests := []struct {
		port     int
		expected error
	}{
		{80, nil},
		{65535, nil},
		{0, errInvalidPort},
		{65536, errInvalidPort},
		{-1, errInvalidPort},
	}

	for _, test := range tests {
		err := validatePort(test.port)
		if !errors.Is(err, test.expected) {
			t.Errorf("validatePort(%q) = %v, expected %v", test.port, err, test.expected)
		}
	}
}

// TestProbeHandler tests the probeHandler function.
func TestProbeHandler(t *testing.T) {
	tests := []struct {
		stateFn  func() bool
		expected int
	}{
		{func() bool { return true }, http.StatusOK},
		{func() bool { return false }, http.StatusServiceUnavailable},
	}

	for _, test := range tests {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler := probeHandler(test.stateFn)
		handler.ServeHTTP(rr, req)

		if rr.Code != test.expected {
			t.Errorf("probeHandler() returned status %v, expected %v", rr.Code, test.expected)
		}
	}
}
