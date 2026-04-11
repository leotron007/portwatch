package daemon_test

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/daemon"
	"github.com/example/portwatch/internal/rules"
)

// captureNotifier records every event it receives.
type captureNotifier struct {
	events []alert.Event
}

func (c *captureNotifier) Notify(ev alert.Event) error {
	c.events = append(c.events, ev)
	return nil
}

func openListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	return l, port
}

func TestDaemon_DetectsOpenedPort(t *testing.T) {
	l, port := openListener(t)
	defer l.Close()

	tmp, err := os.CreateTemp(t.TempDir(), "state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	n := &captureNotifier{}
	cfg := daemon.Config{
		Host:      "127.0.0.1",
		PortRange: fmt.Sprintf("%d-%d", port, port),
		Interval:  20 * time.Millisecond,
		StatePath: tmp.Name(),
		Rules:     []rules.Rule{},
		Notifier:  n,
	}

	d, err := daemon.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	_ = d.Run(ctx)

	if len(n.events) == 0 {
		t.Fatal("expected at least one Opened event, got none")
	}
	if n.events[0].Change != alert.Opened {
		t.Fatalf("expected Opened, got %v", n.events[0].Change)
	}
	if n.events[0].Port != port {
		t.Fatalf("expected port %d, got %d", port, n.events[0].Port)
	}
}

func TestNew_InvalidRange(t *testing.T) {
	_, err := daemon.New(daemon.Config{
		Host:      "127.0.0.1",
		PortRange: "bad-range",
		Interval:  time.Second,
		StatePath: os.DevNull,
		Notifier:  &captureNotifier{},
	})
	if err == nil {
		t.Fatal("expected error for invalid port range")
	}
}
