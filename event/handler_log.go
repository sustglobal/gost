package event

import (
	"fmt"

	"go.uber.org/zap"
)

type logEventPublisher struct {
	logger *zap.Logger
}

func (e *logEventPublisher) HandleEvent(ev *Event) error {
	zfs := make([]zap.Field, len(ev.Fields)+1)
	for i, ef := range ev.Fields {
		zfs[i] = zap.Reflect(fmt.Sprintf("event_field_%s", ef.Key), ef.Value)
	}
	zfs[len(zfs)-1] = zap.String("event_type", string(ev.Type))
	e.logger.Info("event published", zfs...)
	return nil
}

func (e *logEventPublisher) Handles() []EventType {
	return nil
}
