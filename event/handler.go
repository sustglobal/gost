package event

import (
	"context"
)

type EventHandler interface {
	HandleEvent(context.Context, *Event) error
	Handles() []EventType
}
