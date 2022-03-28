package thing

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"path"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/models"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/mqtt"
)

// ThingName the name of the AWS IoT device representation
type ThingName = string

// Thing a structure for working with the AWS IoT device shadows
type Thing struct {
	client    paho.Client
	thingName ThingName
}

// RootPEM is the  Amazon Root CA for IoT Core - Subject to change (but likely not often)
// https://www.amazontrust.com/repository/AmazonRootCA1.pem
const RootPEM = `-----BEGIN CERTIFICATE-----
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

	caPem := []byte(RootPEM)
	certs.AppendCertsFromPEM(caPem)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      certs,
	}

	if err != nil {
		return nil, err
	}

	awsServerURL := fmt.Sprintf("ssl://%s:8883", awsEndpoint)

	mqttOpts := paho.NewClientOptions()
	mqttOpts.AddBroker(awsServerURL)
	mqttOpts.SetMaxReconnectInterval(1 * time.Second)
	mqttOpts.SetClientID(string(thingName))
	mqttOpts.SetTLSConfig(tlsConfig)

	c := paho.NewClient(mqttOpts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Thing{
		client:    c,
		thingName: thingName,
	}, nil
}

// NewThingFromFiles returns a new instance of Thing
func NewThingFromFiles(keyPair models.KeyPair, awsEndpoint string, thingName ThingName) (*Thing, error) {
	client, err := mqtt.MakeMQTTClient(keyPair, awsEndpoint, string(thingName))
	if err != nil {
		return nil, err
	}

	return &Thing{
		client:    client,
		thingName: thingName,
	}, nil
}

// Disconnect terminates the MQTT connection between the client and the AWS server. Recommended to use in defer to avoid
// connection leaks.
func (t *Thing) Disconnect() {
	t.client.Disconnect(1)
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
		func(client paho.Client, msg paho.Message) {
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
