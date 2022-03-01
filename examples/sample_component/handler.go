package sample_component

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/sustglobal/gost/event"
)

var (
	Type_NewDummyEvent = event.EventType("example")
)

func NewDummyHandlerFunc(val string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"CUSTOM_VALUE": %q}`, val)
	}
}

type NewDummyHandler struct {
	*Controller
	*zap.Logger
}

func (h *NewDummyHandler) Handles() []event.EventType {
	return nil
}

func (h *NewDummyHandler) HandleEvent(ctx context.Context, ev *event.Event) error {
	return h.Controller.DoSomething(ctx, ev)
}
