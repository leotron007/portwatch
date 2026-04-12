package notify_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notify"
)

type failChannel struct{ msg string }

func (f *failChannel) Send(_ alert.Event) error { return errors.New(f.msg) }

func TestDispatcher_RegisterAndLen(t *testing.T) {
	d := notify.NewDispatcher()
	if d.Len() != 0 {
		t.Fatalf("expected 0 channels, got %d", d.Len())
	}
	d.Register(notify.NewWriterChannel(&bytes.Buffer{}, ""))
	if d.Len() != 1 {
		t.Fatalf("expected 1 channel, got %d", d.Len())
	}
}

func TestDispatcher_NilChannelIgnored(t *testing.T) {
	d := notify.NewDispatcher()
	d.Register(nil)
	if d.Len() != 0 {
		t.Fatalf("expected 0 channels after nil register, got %d", d.Len())
	}
}

func TestDispatcher_Dispatch_AllChannelsReceiveEvent(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	d := notify.NewDispatcher()
	d.Register(notify.NewWriterChannel(&buf1, "ch1"))
	d.Register(notify.NewWriterChannel(&buf2, "ch2"))

	if err := d.Dispatch(baseEvent); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf1.Len() == 0 || buf2.Len() == 0 {
		t.Error("expected both channels to have received the event")
	}
}

func TestDispatcher_Dispatch_CollectsErrors(t *testing.T) {
	d := notify.NewDispatcher()
	d.Register(&failChannel{msg: "err1"})
	d.Register(&failChannel{msg: "err2"})

	err := d.Dispatch(baseEvent)
	if err == nil {
		t.Fatal("expected error from failing channels")
	}
}

func TestDispatcher_Dispatch_PartialFailure(t *testing.T) {
	var buf bytes.Buffer
	d := notify.NewDispatcher()
	d.Register(notify.NewWriterChannel(&buf, ""))
	d.Register(&failChannel{msg: "boom"})

	err := d.Dispatch(baseEvent)
	if err == nil {
		t.Fatal("expected error from failing channel")
	}
	if buf.Len() == 0 {
		t.Error("successful channel should still have received event")
	}
}
