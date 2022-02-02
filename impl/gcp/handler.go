package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"

	"github.com/sustglobal/gost/component"
	"github.com/sustglobal/gost/event"
)

type MessageHandler interface {
	HandleMessage(ctx context.Context, msg *pubsub.Message) error
}

type gcpPubSubPushRequest struct {
	Message      *pubsub.Message
	Subscription string
}

type PubSubMessageHandler struct {
	MessageHandler
	*zap.Logger
}

func (h *PubSubMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := gcpPubSubPushRequest{
		Message: new(pubsub.Message),
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error("failed decoding PubSub message", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	if err := h.MessageHandler.HandleMessage(r.Context(), req.Message); err != nil {
		h.Logger.Error("failed handling PubSub message", zap.Error(err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

type PubSubMessageEventAdapter struct {
	event.EventHandler
}

func (eh *PubSubMessageEventAdapter) HandleMessage(ctx context.Context, msg *pubsub.Message) error {
	var ev event.Event

	if err := json.Unmarshal(msg.Data, &ev); err != nil {
		return fmt.Errorf("failed unmarshaling PubSub message as event: %v", err)
	}

	return eh.EventHandler.HandleEvent(ctx, &ev)
}

func ListenForPubSubMessages(cmp *component.Component) {
	mh := &PubSubMessageHandler{
		Logger: cmp.Logger,
		MessageHandler: &PubSubMessageEventAdapter{
			EventHandler: cmp.InboundEventRouter,
		},
	}
	cmp.HTTPRouter.Handle("/message", mh).Methods("POST")
}
