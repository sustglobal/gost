package event

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

type fixtureHandler struct {
	types []EventType
	err   error

	events []Event
}

func (h *fixtureHandler) HandleEvent(ctx context.Context, ev *Event) error {
	h.events = append(h.events, *ev)
	return h.err
}

func (h *fixtureHandler) Handles() []EventType {
	return h.types
}

func TestEventRouterUntyped(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	er := NewEventRouter(logger)

	eh := &fixtureHandler{types: nil, err: nil}
	er.Mount(eh)

	ev := Event{Type: EventType("test")}

	if err := er.HandleEvent(context.Background(), &ev); err != nil {
		t.Errorf("received unexpected error: %v", err)
	}

	want := []Event{ev}
	got := eh.events

	if !reflect.DeepEqual(want, got) {
		t.Errorf("events did not route properly: want=%+v got=%+v", want, got)
	}
}

func TestEventRouterTyped(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	er := NewEventRouter(logger)

	eh := &fixtureHandler{types: []EventType{EventType("test2")}, err: nil}
	er.Mount(eh)

	ev1 := Event{Type: EventType("test1")}
	ev2 := Event{Type: EventType("test2")}

	if err := er.HandleEvent(context.Background(), &ev1); err != nil {
		t.Errorf("received unexpected error: %v", err)
	}
	if err := er.HandleEvent(context.Background(), &ev2); err != nil {
		t.Errorf("received unexpected error: %v", err)
	}

	want := []Event{ev2}
	got := eh.events

	if !reflect.DeepEqual(want, got) {
		t.Errorf("events did not route properly: want=%+v got=%+v", want, got)
	}
}

func TestEventRouterMultipleHandlers(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	er := NewEventRouter(logger)

	eh1 := &fixtureHandler{types: []EventType{EventType("test2")}, err: nil}
	er.Mount(eh1)

	eh2 := &fixtureHandler{types: nil, err: nil}
	er.Mount(eh2)

	ev1 := Event{Type: EventType("test1")}
	ev2 := Event{Type: EventType("test2")}

	if err := er.HandleEvent(context.Background(), &ev1); err != nil {
		t.Errorf("received unexpected error: %v", err)
	}
	if err := er.HandleEvent(context.Background(), &ev2); err != nil {
		t.Errorf("received unexpected error: %v", err)
	}

	eh1Want := []Event{ev2}
	eh1Got := eh1.events
	if !reflect.DeepEqual(eh1Want, eh1Got) {
		t.Errorf("events did not route properly to eh1: want=%+v got=%+v", eh1Want, eh1Got)
	}

	eh2Want := []Event{ev1, ev2}
	eh2Got := eh2.events
	if !reflect.DeepEqual(eh2Want, eh2Got) {
		t.Errorf("events did not route properly to eh2: want=%+v got=%+v", eh2Want, eh2Got)
	}

}

func TestEventRouterNoErrorPropagation(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	er := NewEventRouter(logger)

	eh := &fixtureHandler{types: nil, err: errors.New("ignored")}
	er.Mount(eh)

	ev1 := Event{Type: EventType("test1")}
	ev2 := Event{Type: EventType("test2")}

	if err := er.HandleEvent(context.Background(), &ev1); err != nil {
		t.Errorf("received unexpected error: %v", err)
	}

	if err := er.HandleEvent(context.Background(), &ev2); err != nil {
		t.Errorf("received unexpected error: %v", err)
	}

	want := []Event{ev1, ev2}
	got := eh.events

	if !reflect.DeepEqual(want, got) {
		t.Errorf("events did not route properly: want=%+v got=%+v", want, got)
	}
}
