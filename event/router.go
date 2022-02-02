package event

import (
	"context"

	"go.uber.org/zap"
)

func NewEventRouter(logger *zap.Logger) *EventRouter {
	return &EventRouter{
		Logger:          logger,
		typedHandlers:   make(map[EventType][]EventHandler),
		untypedHandlers: make([]EventHandler, 0),
	}
}

type EventRouter struct {
	*zap.Logger
	typedHandlers   map[EventType][]EventHandler
	untypedHandlers []EventHandler
}

func (h *EventRouter) Mount(eh EventHandler) {
	types := eh.Handles()

	if types == nil {
		h.untypedHandlers = append(h.untypedHandlers, eh)
		return
	}

	for _, typ := range types {
		_, ok := h.typedHandlers[typ]
		if !ok {
			h.typedHandlers[typ] = make([]EventHandler, 0)
		}
		h.typedHandlers[typ] = append(h.typedHandlers[typ], eh)
	}
}

//NOTE(bcwaldon): explicitly does NOT handle errors (other than logging) since it is unclear
// what the general behavior should be when a portion of event handlers fail. This may change
// in the future.
func (h *EventRouter) HandleEvent(ctx context.Context, ev *Event) error {
	handlers := make([]EventHandler, 0)
	handlers = append(handlers, h.untypedHandlers...)

	logger := h.Logger.With(zap.String("type", string(ev.Type)))

	typed, ok := h.typedHandlers[ev.Type]
	if ok {
		handlers = append(handlers, typed...)
	}

	if len(handlers) == 0 {
		logger.Debug("no handlers for event")
		return nil
	}

	for _, eh := range handlers {
		if err := eh.HandleEvent(ctx, ev); err != nil {
			logger.Error("event handler failed", zap.Error(err))
		}
	}

	logger.Debug("handled event")

	return nil
}

func (h *EventRouter) Handles() []EventType {
	if len(h.untypedHandlers) > 0 {
		return nil
	}

	types := make([]EventType, 0, len(h.typedHandlers))
	for k, _ := range h.typedHandlers {
		types = append(types, k)
	}
	return types
}
