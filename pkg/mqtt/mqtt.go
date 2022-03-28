package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/models"
)

// MakeMQTTClient creates a new AWS IoT MQTT client.
func MakeMQTTClient(keyPair models.KeyPair, awsEndpoint, clientID string) (mqtt.Client, error) {
	tlsCert, err := tls.LoadX509KeyPair(keyPair.CertificatePath, keyPair.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load the certificates: %v", err)
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
	mqttOpts.SetClientID(clientID)
	mqttOpts.SetTLSConfig(tlsConfig)

	c := mqtt.NewClient(mqttOpts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return c, nil
}
