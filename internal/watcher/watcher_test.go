package watcher_test

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/watcher"
)

func openListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("openListener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func buildWatcher(t *testing.T, minPort, maxPort int, notifier alert.Notifier, buf *bytes.Buffer) *watcher.Watcher {
	t.Helper()

	sc, err := scanner.New(scanner.Config{Host: "127.0.0.1", MinPort: minPort, MaxPort: maxPort})
	if err != nil {
		t.Fatalf("scanner.New: %v", err)
	}

	fl := filter.New(filter.Config{MinPort: minPort, MaxPort: maxPort})

	st, err := state.New(t.TempDir() + "/state.json")
	if err != nil {
		t.Fatalf("state.New: %v", err)
	}

	var rep *reporter.Reporter
	if buf != nil {
		rep = reporter.New(reporter.Config{Writer: buf, Format: "text"})
	}

	w, err := watcher.New(watcher.Config{
		Scanner:  sc,
		Filter:   fl,
		State:    st,
		Rules:    rules.Default(),
		Notifier: notifier,
		Reporter: rep,
		Interval: 50 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("watcher.New: %v", err)
	}
	return w
}

func TestNew_MissingScanner(t *testing.T) {
	_, err := watcher.New(watcher.Config{
		Filter:   filter.New(filter.Config{}),
		State:    &state.State{},
		Notifier: alert.NewLogNotifier(nil),
	})
	if err == nil {
		t.Fatal("expected error for nil scanner")
	}
}

func TestRun_DetectsOpenedPort(t *testing.T) {
	port, closePort := openListener(t)
	defer closePort()

	var logBuf bytes.Buffer
	notifier := alert.NewLogNotifier(&logBuf)

	w := buildWatcher(t, port, port, notifier, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx) // expect context.DeadlineExceeded

	if logBuf.Len() == 0 {
		t.Error("expected at least one alert log line for opened port")
	}
}

func TestRun_ReporterReceivesOutput(t *testing.T) {
	port, closePort := openListener(t)
	defer closePort()

	var repBuf bytes.Buffer
	notifier := alert.NewLogNotifier(nil)

	w := buildWatcher(t, port, port, notifier, &repBuf)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx)

	if repBuf.Len() == 0 {
		t.Error("expected reporter output after scan cycle")
	}
}
