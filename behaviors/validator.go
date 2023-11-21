package behaviors

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/gowool/mediatr"
)

var _ mediatr.PipelineBehavior = (*ValidatorBehavior)(nil)

type Validator interface {
	ValidateCtx(ctx context.Context, i interface{}) error
}

type ValidatorBehavior struct {
	validator Validator
}

func NewValidatorBehavior(validator Validator) ValidatorBehavior {
	return ValidatorBehavior{validator: validator}
}

func (b ValidatorBehavior) Handle(ctx context.Context, request interface{}, next mediatr.RequestHandlerFunc) (interface{}, error) {
	typ := reflect.TypeOf(request)

	logger := slog.Default().WithGroup("mediatr")

	logger.InfoContext(ctx, "validating request", "request_type", typ.String())

	if err := b.validator.ValidateCtx(ctx, request); err != nil {
		logger.WarnContext(ctx, "validation error", "request_type", typ.String(), "request", request, "error", err)

		return nil, err
	}

	return next(ctx)
}
