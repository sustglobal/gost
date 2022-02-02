package gcp

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"

	"github.com/sustglobal/gost/component"
	"github.com/sustglobal/gost/event"
)

func NewPubSubEventPublisher(project, topic string) (*pubsubEventPublisher, error) {
	pubsubClient, err := pubsub.NewClient(context.Background(), project)
	if err != nil {
		return nil, err
	}

	ep := pubsubEventPublisher{
		topic: pubsubClient.Topic(topic),
	}

	return &ep, nil
}

type pubsubEventPublisher struct {
	topic *pubsub.Topic
}

func (p *pubsubEventPublisher) HandleEvent(ctx context.Context, ev *event.Event) error {
	msgData, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	msg := pubsub.Message{
		Data: msgData,
	}

	res := p.topic.Publish(ctx, &msg)

	_, err = res.Get(ctx)
	return err
}

func (p *pubsubEventPublisher) Handles() []event.EventType {
	return nil
}

func PublishEventsToPubSub(cmp *component.Component, gcpProjectID string, gcpPubSubTopic string) error {
	ep, err := NewPubSubEventPublisher(gcpProjectID, gcpPubSubTopic)
	if err != nil {
		return err
	}

	cmp.OutboundEventRouter.Mount(ep)

	return nil
}
