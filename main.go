package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/cloudevents/sdk-go/v01"
	log "github.com/sirupsen/logrus"
)

var projectEnv string
var subscriptionEnv string
var channelEnv string
var namespaceEnv string

func main() {

	//TODO: Make sure all these env vars exist
	projectEnv = os.Getenv("PROJECT")
	subscriptionEnv = os.Getenv("SUBSCRIPTION")
	channelEnv = os.Getenv("CHANNEL")
	namespaceEnv = os.Getenv("NAMESPACE")

	//Setup pubsub client and begin listening for events
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, projectEnv)
	if err != nil {
		log.Error(err)
	}
	subscription := pubsubClient.Subscription(subscriptionEnv)
	subscription.ReceiveSettings.MaxExtension = 10 * time.Second
	log.Info("Listening for events")
	for {
		subscription.Receive(ctx, callback)
	}
}

func callback(ctx context.Context, msg *pubsub.Message) {
	log.Info("Processing msg ID: ", msg.ID)

	//Setup event info
	event := &v01.Event{
		ContentType: "application/json",
		Data:        msg.Data,
		EventID:     msg.ID,
		EventTime:   &msg.PublishTime,
		EventType:   "cloudevent.greet.you",
		Source:      "from-galaxy-far-far-away",
	}

	//Marshal up event JSON and prepare request
	marshaller := v01.NewDefaultHTTPMarshaller()
	req, _ := http.NewRequest("POST", "http://"+channelEnv+"-channel."+namespaceEnv+".svc.cluster.local", nil)
	err := marshaller.ToRequest(req, event)
	if err != nil {
		log.Error(err)
	}

	//Issue POST request, but return before acking the message if there's an error
	_, err = (*http.Client).Do(&http.Client{}, req)
	if err != nil {
		log.Error(err)
		return
	}
	msg.Ack()
}
