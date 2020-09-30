package device

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

// Thing a structure for working with the AWS IoT device shadows
type Thing struct {
	client    mqtt.Client
	thingName ThingName
}

// ThingName the name of the AWS IoT device representation
type ThingName = string

// KeyPair the structure contains the path to the AWS MQTT credentials
type KeyPair struct {
	PrivateKeyPath    string
	CertificatePath   string
	CACertificatePath string
}

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

// Amazon Root CA for IoT Core - Subject to change (but likely not often)
// https://www.amazontrust.com/repository/AmazonRootCA1.pem
const ROOT_PEM = `-----BEGIN CERTIFICATE-----
MIIDQTCCAimgAwIBAgITBmyfz5m/jAo54vB4ikPmljZbyjANBgkqhkiG9w0BAQsF
ADA5MQswCQYDVQQGEwJVUzEPMA0GA1UEChMGQW1hem9uMRkwFwYDVQQDExBBbWF6
b24gUm9vdCBDQSAxMB4XDTE1MDUyNjAwMDAwMFoXDTM4MDExNzAwMDAwMFowOTEL
MAkGA1UEBhMCVVMxDzANBgNVBAoTBkFtYXpvbjEZMBcGA1UEAxMQQW1hem9uIFJv
b3QgQ0EgMTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALJ4gHHKeNXj
ca9HgFB0fW7Y14h29Jlo91ghYPl0hAEvrAIthtOgQ3pOsqTQNroBvo3bSMgHFzZM
9O6II8c+6zf1tRn4SWiw3te5djgdYZ6k/oI2peVKVuRF4fn9tBb6dNqcmzU5L/qw
IFAGbHrQgLKm+a/sRxmPUDgH3KKHOVj4utWp+UhnMJbulHheb4mjUcAwhmahRWa6
VOujw5H5SNz/0egwLX0tdHA114gk957EWW67c4cX8jJGKLhD+rcdqsq08p8kDi1L
93FcXmn/6pUCyziKrlA4b9v7LWIbxcceVOF34GfID5yHI9Y/QCB/IIDEgEw+OyQm
jgSubJrIqg0CAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMC
AYYwHQYDVR0OBBYEFIQYzIU07LwMlJQuCFmcx7IQTgoIMA0GCSqGSIb3DQEBCwUA
A4IBAQCY8jdaQZChGsV2USggNiMOruYou6r4lK5IpDB/G/wkjUu0yKGX9rbxenDI
U5PMCCjjmCXPI6T53iHTfIUJrU6adTrCC2qJeHZERxhlbI1Bjjt/msv0tadQ1wUs
N+gDS63pYaACbvXy8MWy7Vu33PqUXHeeE6V/Uq2V8viTO96LXFvKWlJbYK8U90vv
o/ufQJVtMVT8QtPHRh8jrdkPSHCa2XV4cdFyQzR1bldZwgJcJmApzyMZFo6IQ6XU
5MsI+yMRQ+hDKXJioaldXgjUkK642M4UwtBV8ob2xJNDd2ZhwLnoQdeXeGADbkpy
rqXRfboQnoZsG4q5WTP468SQvvG5
-----END CERTIFICATE-----`

// NewThingFromStrings returns a new instance of Thing
func NewThingFromStrings(cert string, key string, awsEndpoint string, thingName ThingName) (*Thing, error) {
	tlsCert, err := tls.X509KeyPair([]byte(cert), []byte(key))
	certs := x509.NewCertPool()

	caPem := []byte(ROOT_PEM)
	certs.AppendCertsFromPEM(caPem)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      certs,
	}

	if err != nil {
		return nil, err
	}

	awsServerURL := fmt.Sprintf("ssl://%s:8883", awsEndpoint)

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(awsServerURL)
	mqttOpts.SetMaxReconnectInterval(1 * time.Second)
	mqttOpts.SetClientID(string(thingName))
	mqttOpts.SetTLSConfig(tlsConfig)

	c := mqtt.NewClient(mqttOpts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Thing{
		client:    c,
		thingName: thingName,
	}, nil
}

// NewThingFromFiles returns a new instance of Thing
func NewThingFromFiles(keyPair KeyPair, awsEndpoint string, thingName ThingName) (*Thing, error) {
	tlsCert, err := tls.LoadX509KeyPair(keyPair.CertificatePath, keyPair.PrivateKeyPath)
	if err != nil {
		return nil ,fmt.Errorf("failed to load the certificates: %v", err)
	}

	certs := x509.NewCertPool()

	caPem, err := ioutil.ReadFile(keyPair.CACertificatePath)
	if err != nil {
		return nil, err
	}

	certs.AppendCertsFromPEM(caPem)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      certs,
	}

	if err != nil {
		return nil, err
	}

	awsServerURL := fmt.Sprintf("ssl://%s:8883", awsEndpoint)

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(awsServerURL)
	mqttOpts.SetMaxReconnectInterval(1 * time.Second)
	mqttOpts.SetClientID(string(thingName))
	mqttOpts.SetTLSConfig(tlsConfig)

	c := mqtt.NewClient(mqttOpts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Thing{
		client:    c,
		thingName: thingName,
	}, nil
}

// Disconnect terminates the MQTT connection between the client and the AWS server. Recommended to use in defer to avoid
// connection leaks.
func (t *Thing) Disconnect() {
	t.client.Disconnect(1)
}

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

// PublishToCustomTopic publishes an async message to the custom topic.
func (t *Thing) PublishToCustomTopic(payload Payload, topic string) error {
	token := t.client.Publish(
		topic,
		0,
		false,
		[]byte(payload),
	)
	token.Wait()
	return token.Error()
}

// SubscribeForCustomTopic subscribes for the custom topic and returns the channel with the topic messages.
func (t *Thing) SubscribeForCustomTopic(topic string) (chan Payload, error) {
	payloadChan := make(chan Payload)

	if token := t.client.Subscribe(
		topic,
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			payloadChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return payloadChan, nil
}

// UnsubscribeFromCustomTopic terminates the subscription to the custom topic.
// The specified topic argument will be prepended by a prefix "$aws/things/<thing_name>"
func (t Thing) UnsubscribeFromCustomTopic(topic string) error {
	return t.unsubscribe(path.Join("$aws/things", t.thingName, topic))
}

// unsubscribe terminates the MQTT subscription for the provided tokens
func (t Thing) unsubscribe(topics ...string) error {
	token := t.client.Unsubscribe(topics...)
	token.Wait()
	return token.Error()
}
