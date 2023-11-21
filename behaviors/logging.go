package behaviors

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/gowool/mediatr"
)

var _ mediatr.PipelineBehavior = (*LoggingBehavior)(nil)

type LoggingBehavior struct{}

func NewLoggingBehavior() LoggingBehavior {
	return LoggingBehavior{}
}

func (b LoggingBehavior) Handle(ctx context.Context, request interface{}, next mediatr.RequestHandlerFunc) (interface{}, error) {
	typ := reflect.TypeOf(request)

	logger := slog.Default().WithGroup("mediatr")

	logger.InfoContext(ctx, "handling request", "request_type", typ.String(), "request", request)

	response, err := next(ctx)

	logger.InfoContext(ctx, "command handled - response", "request_type", typ.String(), "response", response)

	return response, err
}
