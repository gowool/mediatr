package mediatr

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/gowool/mediatr/dict"
)

var (
	ErrBehaviorConflict = errors.New("behavior already was registered")
	ErrHandlerConflict  = errors.New("request handler already was registered")
	ErrHandlerNotFound  = errors.New("request handler not found")
)

var (
	requestHandlers   = dict.NewDict[reflect.Type, any]()
	pipelineBehaviors = dict.NewDict[reflect.Type, PipelineBehavior]()
)

type Unit struct{}

type RequestHandlerFunc func(ctx context.Context) (interface{}, error)

type PipelineBehavior interface {
	Handle(ctx context.Context, request interface{}, next RequestHandlerFunc) (interface{}, error)
}

type RequestHandler[TRequest any, TResponse any] interface {
	Handle(ctx context.Context, request TRequest) (TResponse, error)
}

type RequestHandlerFactory[TRequest any, TResponse any] func(ctx context.Context) (RequestHandler[TRequest, TResponse], error)

func ClearRequestHandlers() {
	requestHandlers.Clear()
}

func ClearPipelineBehaviors() {
	pipelineBehaviors.Clear()
}

func registerRequestHandler[TRequest any](handler interface{}) error {
	var request TRequest
	typ := reflect.TypeOf(request)

	if requestHandlers.Has(typ) {
		return fmt.Errorf("`%s` %w", typ.String(), ErrHandlerConflict)
	}

	requestHandlers.Set(typ, handler)

	return nil
}

func RegisterRequestHandler[TRequest any, TResponse any](handler RequestHandler[TRequest, TResponse]) error {
	return registerRequestHandler[TRequest](handler)
}

func RegisterRequestHandlerFactory[TRequest any, TResponse any](factory RequestHandlerFactory[TRequest, TResponse]) error {
	return registerRequestHandler[TRequest](factory)
}

func RegisterPipelineBehaviors(behaviors ...PipelineBehavior) error {
	for _, behavior := range behaviors {
		typ := reflect.TypeOf(behavior)

		if pipelineBehaviors.Has(typ) {
			return fmt.Errorf("`%s` %w", typ.String(), ErrBehaviorConflict)
		}

		pipelineBehaviors.Set(typ, behavior)
	}
	return nil
}

func Send[TRequest any](ctx context.Context, request TRequest) error {
	_, err := SendR[TRequest, Unit](ctx, request)
	return err
}

func SendR[TRequest any, TResponse any](ctx context.Context, request TRequest) (TResponse, error) {
	var response TResponse

	hi, ok := requestHandlers.Get(reflect.TypeOf(request))

	if !ok {
		return response, fmt.Errorf("`%T` %w", request, ErrHandlerNotFound)
	}

	handler, err := toRequestHandler[TRequest, TResponse](ctx, hi)
	if err != nil {
		return response, err
	}

	var handlerFunc RequestHandlerFunc = func(ctx context.Context) (interface{}, error) {
		return handler.Handle(ctx, request)
	}

	behaviors := pipelineBehaviors.Values()
	for i := len(behaviors) - 1; i >= 0; i-- {
		next := handlerFunc
		pipe := behaviors[i]

		handlerFunc = func(ctx context.Context) (interface{}, error) {
			return pipe.Handle(ctx, request, next)
		}
	}

	value, err := handlerFunc(ctx)
	if err != nil {
		return response, err
	}

	return value.(TResponse), nil
}

func toRequestHandler[TRequest any, TResponse any](ctx context.Context, i interface{}) (RequestHandler[TRequest, TResponse], error) {
	switch handler := i.(type) {
	case RequestHandler[TRequest, TResponse]:
		return handler, nil
	case RequestHandlerFactory[TRequest, TResponse]:
		return handler(ctx)
	default:
		var request TRequest
		panic(fmt.Errorf("handler for request %T is not a handler", request))
	}
}
