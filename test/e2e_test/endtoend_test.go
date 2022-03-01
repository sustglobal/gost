package e2e_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sustglobal/gost/component"
	"github.com/sustglobal/gost/examples/sample_component"
	"github.com/sustglobal/gost/impl/gcp"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type ServiceConfig struct {
	CustomValue string `env:"CUSTOM_VALUE"`
}

type PubSubConfig struct {
	GCPProjectID        string
	GCPPubSubTopic      string
	GCPSubscriptionName string
	context             context.Context
	options             []option.ClientOption
}

func TestE2E(t *testing.T) {
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	pubsub_cfg := PubSubConfig{GCPProjectID: "Test", GCPPubSubTopic: "TestTopic", GCPSubscriptionName: "TestSub", context: ctx, options: []option.ClientOption{option.WithoutAuthentication()}}
	setup_client(pubsub_cfg)

	cmp, err := component.NewFromEnv()
	if err != nil {
		panic(err)
	}

	var service_cfg ServiceConfig

	if err := component.LoadFromEnv(&service_cfg); err != nil {
		panic(err)
	}

	cont := &sample_component.Controller{
		Data:   "TEST!",
		Logger: cmp.Logger,
	}

	cmp.InboundEventRouter.Mount(
		&sample_component.NewDummyHandler{
			Controller: cont,
			Logger:     cmp.Logger,
		},
	)

	cmp.OutboundEventRouter.Mount(
		&sample_component.NewDummyHandler{
			Controller: cont,
			Logger:     cmp.Logger,
		},
	)

	gcp.PublishEventsToPubSub(cmp, pubsub_cfg.GCPProjectID, pubsub_cfg.GCPPubSubTopic)
	gcp.ListenForPubSubMessages(cmp)

	cmp.Run()
}

func setup_client(pubsub_cfg PubSubConfig) {
	client, _ := pubsub.NewClient(pubsub_cfg.context, pubsub_cfg.GCPProjectID)

	topic, _ := client.CreateTopic(pubsub_cfg.context, pubsub_cfg.GCPPubSubTopic)
	client.CreateSubscription(pubsub_cfg.context, pubsub_cfg.GCPSubscriptionName, pubsub.SubscriptionConfig{Topic: topic})
}
