package tmevents

import (
	"net/http"
	"time"

	"github.com/cloudevents/sdk-go/v01"
)

// EventInfo holds info about the event that occurred
type EventInfo struct {
	EventData   []byte
	EventID     string
	EventTime   time.Time
	EventType   string
	EventSource string
}

// TMEventInterface is here to allow each type of message watcher
// to implement however it'll receive message data
type TMEventInterface interface {
	ReceiveMsg()
}

// PushEvent pushes an event to a kubernetes service
func PushEvent(ev *EventInfo, desiredURL string) error {

	//Setup event info
	event := &v01.Event{
		ContentType: "application/json",
		Data:        ev.EventData,
		EventID:     ev.EventID,
		EventTime:   &ev.EventTime,
		EventType:   ev.EventType,
		Source:      ev.EventSource,
	}

	//Marshal up event JSON and prepare request
	marshaller := v01.NewDefaultHTTPMarshaller()
	req, _ := http.NewRequest("POST", desiredURL, nil)
	err := marshaller.ToRequest(req, event)
	if err != nil {
		return err
	}

	//Issue POST request
	_, err = (*http.Client).Do(&http.Client{}, req)
	if err != nil {
		return err
	}

	return nil
}
