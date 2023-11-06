package behaviors

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/gowool/mediatr"
)

var _ mediatr.PipelineBehavior = (*LoggingBehavior)(nil)

type LoggingBehavior struct {
	logger *slog.Logger
}

func NewLoggingBehavior(logger *slog.Logger) LoggingBehavior {
	return LoggingBehavior{logger: logger}
}

func (b LoggingBehavior) Handle(ctx context.Context, request interface{}, next mediatr.RequestHandlerFunc) (interface{}, error) {
	typ := reflect.TypeOf(request)

	b.logger.InfoContext(ctx, "handling request", "request_type", typ.String(), "request", request)

	response, err := next(ctx)

	b.logger.InfoContext(ctx, "command handled - response", "request_type", typ.String(), "response", response)

	return response, err
}
