package handler

import (
	"context"
	"encoding/json"

	"log/slog"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	slogcommon "github.com/samber/slog-common"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// Whether to log as structured (default: false)
	StructLog bool

	// Whether to pretty-print the log output (default: false)
	Pretty bool

	// optional: zerolog logger (default: zerolog.Logger)
	Logger *zerolog.Logger

	// optional: customize json payload builder
	Converter Converter

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
}

func (o Option) NewZerologHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Logger == nil {
		// should be selected lazily ?
		o.Logger = &log.Logger
	}

	return &ZerologHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*ZerologHandler)(nil)

type ZerologHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *ZerologHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *ZerologHandler) Handle(ctx context.Context, record slog.Record) error {
	converter := DefaultConverter
	if h.option.Converter != nil {
		converter = h.option.Converter
	}

	level := LogLevels[record.Level]
	args := converter(h.option.AddSource, h.option.ReplaceAttr, h.attrs, h.groups, &record)

	if h.option.StructLog || len(args) == 0 {
		h.option.Logger.
			WithLevel(level).
			Ctx(ctx).
			Time(zerolog.TimestampFieldName, record.Time.UTC()).
			Fields(args).
			Msg(record.Message)
	} else {
		var msg []byte
		if h.option.Pretty {
			msg, _ = json.MarshalIndent(args, "", "  ")
		} else {
			msg, _ = json.Marshal(args)
		}

		h.option.Logger.
			WithLevel(level).
			Ctx(ctx).
			Time(zerolog.TimestampFieldName, record.Time.UTC()).
			Msgf("%s %s", record.Message, msg)
	}

	return nil
}

func (h *ZerologHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ZerologHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *ZerologHandler) WithGroup(name string) slog.Handler {
	return &ZerologHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
