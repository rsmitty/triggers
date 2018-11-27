package main

import (
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rsmitty/triggers/tmevents"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
)

var queueEnv string
var channelEnv string
var namespaceEnv string
var awsaccessKeyEnv string
var awsSecretKeyEnv string
var awsRegionEnv string
var awsCredsFile string

type sqsMsg struct{}

func main() {

	//TODO: Make sure all these env vars exist
	queueEnv = os.Getenv("QUEUE")
	channelEnv = os.Getenv("CHANNEL")
	namespaceEnv = os.Getenv("NAMESPACE")
	awsRegionEnv = os.Getenv("AWS_REGION")
	awsCredsFile = os.Getenv("AWS_CREDS")

	//Create client for SQS and start polling for messages on the queue
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegionEnv),
		Credentials: credentials.NewSharedCredentials(awsCredsFile, "default"),
		MaxRetries:  aws.Int(5),
	})
	if err != nil {
		log.Fatal(err)
	}
	sqsClient := sqs.New(sess)
	m := sqsMsg{}
	m.ReceiveMsg(sqsClient)

}

//ReceiveMsg implements the receive interface for sqs
func (sqsMsg) ReceiveMsg(sqsClient *sqs.SQS) {
	//Look for new messages every 5 seconds
	for range time.Tick(5 * time.Second) {
		msg, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: aws.StringSlice([]string{"All"}),
			QueueUrl:       &queueEnv,
		})
		if err != nil {
			log.Info(err)
			continue
		}
		//Only push if there are messages on the queue
		if len(msg.Messages) > 0 {
			log.Info("Processing message with ID: ", aws.StringValue(msg.Messages[0].MessageId))
			log.Info(msg.Messages[0])
			//Parse out timestamp
			msgAttributes := aws.StringValueMap(msg.Messages[0].Attributes)
			timeInt, err := strconv.ParseInt(msgAttributes["SentTimestamp"], 10, 64)
			if err != nil {
				log.Info(err)
				continue
			}
			timeSent := time.Unix(timeInt, 0)

			//Craft event info and push it
			eventInfo := tmevents.EventInfo{
				EventData:   []byte(aws.StringValue(msg.Messages[0].Body)),
				EventID:     aws.StringValue(msg.Messages[0].MessageId),
				EventTime:   timeSent,
				EventType:   "cloudevent.greet.you",
				EventSource: "sqs",
			}
			url := "http://" + channelEnv + "-channel." + namespaceEnv + ".svc.cluster.local"
			err = tmevents.PushEvent(&eventInfo, url)
			if err != nil {
				log.Error(err)
				continue
			}

			//Delete message from queue if we pushed successfully
			err = deleteMessage(sqsClient, msg)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

//Deletes message from sqs queue
func deleteMessage(sqsClient *sqs.SQS, msg *sqs.ReceiveMessageOutput) error {
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueEnv),
		ReceiptHandle: msg.Messages[0].ReceiptHandle,
	}
	_, err := sqsClient.DeleteMessage(deleteParams)
	if err != nil {
		return err
	}
	return nil
}
