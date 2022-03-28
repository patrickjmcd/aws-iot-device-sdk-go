package thing

import (
	"errors"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Shadow device shadow data
type Shadow []byte

// Payload device data
type Payload []byte

// String converts the Shadow to string
func (s Shadow) String() string {
	return string(s)
}

// ShadowError represents the model for handling the errors occurred during updating the device shadow
type ShadowError = Shadow

// GetThingShadow returns the current thing shadow
func (t *Thing) GetThingShadow() (Shadow, error) {
	shadowChan := make(chan Shadow)
	errChan := make(chan error)

	defer t.unsubscribe(
		fmt.Sprintf("$aws/things/%s/shadow/get/accepted", t.thingName),
		fmt.Sprintf("$aws/things/%s/shadow/get/rejected", t.thingName),
	)

	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/shadow/get/accepted", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			shadowChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/shadow/get/rejected", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			errChan <- errors.New(string(msg.Payload()))
		},
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	if token := t.client.Publish(
		fmt.Sprintf("$aws/things/%s/shadow/get", t.thingName),
		0,
		false,
		[]byte("{}"),
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	for {
		select {
		case s, ok := <-shadowChan:
			if !ok {
				return nil, errors.New("failed to read from shadow channel")
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

// UpdateThingShadow publishes an async message with new thing shadow
func (t *Thing) UpdateThingShadow(payload Shadow) error {
	token := t.client.Publish(fmt.Sprintf("$aws/things/%s/shadow/update", t.thingName), 1, false, []byte(payload))
	token.Wait()
	return token.Error()
}

// SubscribeForThingShadowChanges subscribes for the device shadow update topic and returns two channels: shadow and shadow error.
// The shadow channel will handle all accepted device shadow updates. The shadow error channel will handle all rejected device
// shadow updates
func (t *Thing) SubscribeForThingShadowChanges() (chan Shadow, chan ShadowError, error) {
	shadowChan := make(chan Shadow)
	shadowErrChan := make(chan ShadowError)

	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/shadow/update/accepted", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			shadowChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return nil, nil, token.Error()
	}

	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/shadow/update/rejected", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			shadowErrChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return nil, nil, token.Error()
	}

	return shadowChan, shadowErrChan, nil
}

// UpdateThingShadowDocument publishes an async message with new thing shadow document
func (t *Thing) UpdateThingShadowDocument(payload Shadow) error {
	token := t.client.Publish(fmt.Sprintf("$aws/things/%s/shadow/update/documents", t.thingName), 0, false, []byte(payload))
	token.Wait()
	return token.Error()
}

// DeleteThingShadow publishes a message to remove the device's shadow and waits for the result. In case shadow delete was
// rejected the method will return error
func (t *Thing) DeleteThingShadow() error {
	shadowChan := make(chan Shadow)
	errChan := make(chan error)

	defer t.unsubscribe(
		fmt.Sprintf("$aws/things/%s/shadow/delete/accepted", t.thingName),
		fmt.Sprintf("$aws/things/%s/shadow/delete/rejected", t.thingName),
	)

	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/shadow/delete/accepted", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			shadowChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := t.client.Subscribe(
		fmt.Sprintf("$aws/things/%s/shadow/delete/rejected", t.thingName),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			errChan <- errors.New(string(msg.Payload()))
		},
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := t.client.Publish(
		fmt.Sprintf("$aws/things/%s/shadow/delete", t.thingName),
		0,
		false,
		[]byte("{}"),
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for {
		select {
		case _, ok := <-shadowChan:
			if !ok {
				return errors.New("failed to read from shadow channel")
			}
			return nil
		case err, ok := <-errChan:
			if !ok {
				return errors.New("failed to read from error channel")
			}
			return err
		}
	}
}
