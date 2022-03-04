package sample_component

import (
	"context"

	"github.com/sustglobal/gost/event"

	"go.uber.org/zap"
)

type Controller struct {
	EventChan chan<- *event.Event
	*zap.Logger
}

func (c *Controller) DoSomething(ctx context.Context, ev *event.Event) error {
	c.Logger.Info("controller.DoSomething received event", zap.Reflect("event", ev))
	c.EventChan <- ev
	return nil
}
