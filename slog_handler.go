package gologger_nazalog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/q191201771/naza/pkg/nazalog"
)

type SlogHandler struct {
	l      nazalog.Logger
	opts   slog.HandlerOptions
	attrs  []slog.Attr
	groups []string
}

func NewSlogHandler(l nazalog.Logger, opts *slog.HandlerOptions) *SlogHandler {
	if l == nil {
		l = nazalog.GetGlobalLogger()
	}
	h := &SlogHandler{l: l}
	if opts != nil {
		h.opts = *opts
	}
	return h
}

func NewSlogLogger(l nazalog.Logger, opts *slog.HandlerOptions) *slog.Logger {
	return slog.New(NewSlogHandler(l, opts))
}

func (h *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	min := slog.LevelInfo
	if h.opts.Level != nil {
		min = h.opts.Level.Level()
	}
	if level < min {
		return false
	}
	if h.l == nil {
		return true
	}
	return h.l.GetOption().Level <= slogLevelToNazalog(level)
}

func (h *SlogHandler) Handle(_ context.Context, r slog.Record) error {
	if h.l == nil {
		return nil
	}

	var buf bytes.Buffer
	buf.WriteString(r.Message)

	if h.opts.AddSource && r.PC != 0 {
		file, line := sourceFromPC(r.PC)
		if file != "" && line > 0 {
			appendKV(&buf, "source", file+":"+strconv.Itoa(line))
		}
	}

	appendAttrs(&buf, h.groups, h.attrs, h.opts.ReplaceAttr)
	r.Attrs(func(a slog.Attr) bool {
		appendAttr(&buf, h.groups, a, h.opts.ReplaceAttr)
		return true
	})

	h.l.Out(slogLevelToNazalog(r.Level), callDepthForOut(), buf.String())
	return nil
}

func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	cp := h.clone()
	cp.attrs = append(cp.attrs, attrs...)
	return cp
}

func (h *SlogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	cp := h.clone()
	cp.groups = append(cp.groups, name)
	return cp
}

func (h *SlogHandler) clone() *SlogHandler {
	cp := *h
	if len(h.attrs) > 0 {
		cp.attrs = append([]slog.Attr(nil), h.attrs...)
	}
	if len(h.groups) > 0 {
		cp.groups = append([]string(nil), h.groups...)
	}
	return &cp
}

func slogLevelToNazalog(level slog.Level) nazalog.Level {
	if level < slog.LevelDebug {
		return nazalog.LevelTrace
	}
	if level < slog.LevelInfo {
		return nazalog.LevelDebug
	}
	if level < slog.LevelWarn {
		return nazalog.LevelInfo
	}
	if level < slog.LevelError {
		return nazalog.LevelWarn
	}
	return nazalog.LevelError
}

func sourceFromPC(pc uintptr) (string, int) {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	return f.File, f.Line
}

func callDepthForOut() int {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	depthFromHandle := 0
	for {
		f, more := frames.Next()
		if depthFromHandle == 0 {
			depthFromHandle++
			if !more {
				break
			}
			continue
		}

		if isInternalSlogFrame(f.Function) || isInternalNazalogFrame(f.Function) || isRuntimeFrame(f.Function) {
			depthFromHandle++
			if !more {
				break
			}
			continue
		}

		return depthFromHandle + 1
	}

	return 3
}

func isInternalSlogFrame(fn string) bool {
	return strings.Contains(fn, "log/slog.")
}

func isInternalNazalogFrame(fn string) bool {
	return strings.Contains(fn, "github.com/kordar/gologger_nazalog.") || strings.Contains(fn, ".(*SlogHandler).")
}

func isRuntimeFrame(fn string) bool {
	return strings.HasPrefix(fn, "runtime.")
}

func appendAttrs(buf *bytes.Buffer, groups []string, attrs []slog.Attr, replace func(groups []string, a slog.Attr) slog.Attr) {
	for _, a := range attrs {
		appendAttr(buf, groups, a, replace)
	}
}

func appendAttr(buf *bytes.Buffer, groups []string, a slog.Attr, replace func(groups []string, a slog.Attr) slog.Attr) {
	a.Value = a.Value.Resolve()
	if replace != nil {
		a = replace(groups, a)
	}
	if a.Equal(slog.Attr{}) {
		return
	}
	if a.Key == "" && a.Value.Kind() != slog.KindGroup {
		return
	}

	if a.Value.Kind() == slog.KindGroup {
		ng := groups
		if a.Key != "" {
			ng = append(append([]string(nil), groups...), a.Key)
		}
		for _, ca := range a.Value.Group() {
			appendAttr(buf, ng, ca, replace)
		}
		return
	}

	appendKV(buf, keyWithGroups(groups, a.Key), formatValue(a.Value))
}

func appendKV(buf *bytes.Buffer, key string, value string) {
	if key == "" {
		return
	}
	buf.WriteByte(' ')
	buf.WriteString(key)
	buf.WriteByte('=')
	buf.WriteString(value)
}

func keyWithGroups(groups []string, key string) string {
	if len(groups) == 0 {
		return key
	}
	parts := make([]string, 0, len(groups)+1)
	for _, g := range groups {
		if g == "" {
			continue
		}
		parts = append(parts, g)
	}
	parts = append(parts, key)
	return strings.Join(parts, ".")
}

func formatValue(v slog.Value) string {
	switch v.Kind() {
	case slog.KindString:
		return quoteIfNeeded(v.String())
	case slog.KindInt64:
		return strconv.FormatInt(v.Int64(), 10)
	case slog.KindUint64:
		return strconv.FormatUint(v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.FormatFloat(v.Float64(), 'f', -1, 64)
	case slog.KindBool:
		return strconv.FormatBool(v.Bool())
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().Format(time.RFC3339Nano)
	default:
		return quoteIfNeeded(fmt.Sprint(v.Any()))
	}
}

func quoteIfNeeded(s string) string {
	if s == "" {
		return `""`
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r', '=':
			return strconv.Quote(s)
		}
	}
	return s
}
