package gologger_nazalog_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/kordar/gologger_nazalog"
	"github.com/q191201771/naza/pkg/nazalog"
)

func TestSlogHandler_BasicFields(t *testing.T) {
	var got bytes.Buffer

	l, err := nazalog.New(func(o *nazalog.Option) {
		o.Level = nazalog.LevelTrace
		o.IsToStdout = false
		o.Filename = ""
		o.TimestampFlag = false
		o.LevelFlag = false
		o.ShortFileFlag = false
		o.HookBackendOutFn = func(level nazalog.Level, line []byte) {
			got.Write(line)
		}
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	sl := gologger_nazalog.NewSlogLogger(l, &slog.HandlerOptions{Level: slog.LevelDebug})
	sl.Info("m", "k", "v")

	s := got.String()
	if !strings.Contains(s, "m") {
		t.Fatalf("expected message in output, got: %q", s)
	}
	if !strings.Contains(s, "k=v") {
		t.Fatalf("expected field in output, got: %q", s)
	}
}

func TestSlogHandler_GroupAndAttrs(t *testing.T) {
	var got bytes.Buffer

	l, err := nazalog.New(func(o *nazalog.Option) {
		o.Level = nazalog.LevelTrace
		o.IsToStdout = false
		o.Filename = ""
		o.TimestampFlag = false
		o.LevelFlag = false
		o.ShortFileFlag = false
		o.HookBackendOutFn = func(level nazalog.Level, line []byte) {
			got.Write(line)
		}
	})
	if err != nil {
		t.Fatalf("new logger: %v", err)
	}

	h := gologger_nazalog.NewSlogHandler(l, &slog.HandlerOptions{Level: slog.LevelDebug})
	sl := slog.New(h.WithGroup("g").WithAttrs([]slog.Attr{slog.String("a", "b")}))

	sl.Info("m", "k", "v")

	s := got.String()
	if !strings.Contains(s, "g.k=v") {
		t.Fatalf("expected grouped key in output, got: %q", s)
	}
	if !strings.Contains(s, "g.a=b") {
		t.Fatalf("expected grouped attr in output, got: %q", s)
	}
}
