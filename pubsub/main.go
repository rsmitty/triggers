package main

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/rsmitty/tmevents"

	log "github.com/sirupsen/logrus"
)

var projectEnv string
var subscriptionEnv string
var channelEnv string
var namespaceEnv string

type pubsubMsg struct{}

func main() {

	//TODO: Make sure all these env vars exist
	projectEnv = os.Getenv("PROJECT")
	subscriptionEnv = os.Getenv("SUBSCRIPTION")
	channelEnv = os.Getenv("CHANNEL")
	namespaceEnv = os.Getenv("NAMESPACE")

	m := pubsubMsg{}
	m.ReceiveMsg()

}

//ReceiveMsg implements the the receive interface for pubsub
func (m pubsubMsg) ReceiveMsg() {
	//Setup pubsub client and begin listening for events
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, projectEnv)
	if err != nil {
		log.Error(err)
	}

	//Listen forever for messages from our subscription
	subscription := pubsubClient.Subscription(subscriptionEnv)
	subscription.ReceiveSettings.MaxExtension = 10 * time.Second
	log.Info("Listening for events")
	for {
		subscription.Receive(ctx, callback)
	}
}

//Pulls msg info, creates a tmevents client, and pushes the event.
//Returns before acking msg if there's an error.
func callback(ctx context.Context, msg *pubsub.Message) {
	log.Info("Processing msg ID: ", msg.ID)

	eventInfo := tmevents.EventInfo{
		EventData:   msg.Data,
		EventID:     msg.ID,
		EventTime:   msg.PublishTime,
		EventType:   "cloudevent.greet.you",
		EventSource: "pubsub",
	}
	url := "http://" + channelEnv + "-channel." + namespaceEnv + ".svc.cluster.local"
	err := tmevents.PushEvent(&eventInfo, url)
	if err != nil {
		log.Error(err)
		return
	}
	msg.Ack()
}
