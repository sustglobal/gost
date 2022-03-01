package sample_component

import (
	"context"
	"fmt"
	"log"

	"github.com/sustglobal/gost/event"

	"go.uber.org/zap"
)

type Controller struct {
	Data string
	*zap.Logger
}

func (c *Controller) DoSomething(ctx context.Context, ev *event.Event) error {
	// Doesn't actually do anything...
	log.Println("I did a thing!")
	if ev == nil {
		return fmt.Errorf("something is wrong")
	}
	return nil
}
