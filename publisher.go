package mediatr

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/gowool/mediatr/dict"
	"github.com/gowool/mediatr/list"
)

var ErrStop = errors.New("stop")

var notificationHandlers = dict.NewDict[reflect.Type, *list.List[any]]()

type NotificationHandler[TNotification any] interface {
	Handle(ctx context.Context, notification TNotification) error
}

type NotificationHandlerFactory[TNotification any] func(ctx context.Context) (NotificationHandler[TNotification], error)

func ClearNotificationHandlers() {
	notificationHandlers.Loop(func(_ reflect.Type, u *list.List[any]) bool {
		u.Clear()
		return true
	})
	notificationHandlers.Clear()
}

func registerNotificationHandler[TNotification any](handler interface{}) error {
	var notification TNotification
	typ := reflect.TypeOf(notification)

	if handlers, ok := notificationHandlers.Get(typ); ok {
		handlers.Add(handler)
	} else {
		handlers = list.NewList(handler)
		notificationHandlers.Set(typ, handlers)
	}
	return nil
}

func RegisterNotificationHandlers[TNotification any](handlers ...NotificationHandler[TNotification]) error {
	for _, handler := range handlers {
		if err := registerNotificationHandler[TNotification](handler); err != nil {
			return err
		}
	}
	return nil
}

func RegisterNotificationHandlerFactories[TNotification any](factories ...NotificationHandlerFactory[TNotification]) error {
	for _, factory := range factories {
		if err := registerNotificationHandler[TNotification](factory); err != nil {
			return err
		}
	}
	return nil
}

func Publish[TNotification any](ctx context.Context, notification TNotification, oneOff ...interface{}) (err error) {
	handlers, ok := notificationHandlers.Get(reflect.TypeOf(notification))
	if !ok && len(oneOff) == 0 {
		// notification strategy should have zero or more handlers, so it should run without any error if we can't find a corresponding handler
		return
	}

	items := make([]interface{}, 0, handlers.Len()+len(oneOff))
	items = append(items, handlers.Items()...)
	items = append(items, oneOff...)

	for _, item := range items {
		var handler NotificationHandler[TNotification]
		if handler, err = toNotificationHandler[TNotification](ctx, item); err != nil {
			return
		}

		if err1 := handler.Handle(ctx, notification); err1 != nil {
			if !errors.Is(err1, ErrStop) {
				err = err1
			}
			return
		}
	}

	return
}

func toNotificationHandler[TNotification any](ctx context.Context, i interface{}) (NotificationHandler[TNotification], error) {
	switch handler := i.(type) {
	case NotificationHandler[TNotification]:
		return handler, nil
	case NotificationHandlerFactory[TNotification]:
		return handler(ctx)
	default:
		var notification TNotification
		panic(fmt.Errorf("handler for notification %T is not a handler", notification))
	}
}
