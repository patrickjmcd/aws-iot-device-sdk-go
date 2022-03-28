package thing

import (
	"errors"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ListenForJobs is a helper function that subscribes to the topic responsible for notifying on IoT Core Jobs
func (t *Thing) ListenForJobs() (chan Payload, error) {
	// return t.SubscribeForCustomTopic(fmt.Sprintf("$aws/things/%s/jobs/notify", t.thingName))
	jobsChan := make(chan Payload)
	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/jobs/notify", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			jobsChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/jobs/notify-next", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			jobsChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/jobs/get/accepted", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			jobsChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/jobs/get/rejected", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			jobsChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return jobsChan, nil
}

func (t *Thing) GetNextJob() (Payload, error) {
	jobsChan := make(chan Payload)
	errChan := make(chan error)

	defer t.unsubscribe(
		fmt.Sprintf("$aws/things/%s/jobs/get/#", t.thingName),
		fmt.Sprintf("$aws/things/%s/jobs/next", t.thingName),
		fmt.Sprintf("$aws/things/%s/jobs/next-notify", t.thingName),
	)

	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/jobs/next", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			jobsChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/jobs/next-notify", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			jobsChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	if token := t.client.Publish(
		fmt.Sprintf("$aws/things/%s/jobs/get", t.thingName),
		0,
		false,
		[]byte(fmt.Sprintf("{%s: %s, %s: %s}", "clientToken", t.thingName, "jobId", "$next")),
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	for {
		select {
		case s, ok := <-jobsChan:
			if !ok {
				return nil, errors.New("failed to read from jobs channel")
			}
			return s, nil
		case err, ok := <-errChan:
			if !ok {
				return nil, errors.New("failed to read from error channel")
			}
			return nil, err
		}
	}
}

//
func (t *Thing) UnsubscribeFromJobs() error {
	err := t.unsubscribe(fmt.Sprintf("$aws/things/%s/jobs/notify", t.thingName))
	err = t.unsubscribe(fmt.Sprintf("$aws/things/%s/jobs/notify-next", t.thingName))
	err = t.unsubscribe(fmt.Sprintf("$aws/things/%s/jobs/get/accepted", t.thingName))
	err = t.unsubscribe(fmt.Sprintf("$aws/things/%s/jobs/get/rejected", t.thingName))
	return err
}
