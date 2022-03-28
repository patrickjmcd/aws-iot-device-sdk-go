package tunnel

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/models"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/mqtt"
)

type tunnelPayload struct {
	ClientAccessToken string   `json:"clientAccessToken"`
	ClientMode        string   `json:"clientMode"`
	Region            string   `json:"region"`
	Services          []string `json:"services"`
}

// ListenForTunnel listens on the MQTT Tunnel Topic and sets up the tunnel once a notify message is received
func ListenForTunnel(thingName string, keypair models.KeyPair, endpoint string) error {

	_, err := exec.LookPath("localproxy")
	if err != nil {
		return fmt.Errorf("localproxy not found in path")
	}

	notifyChan := make(chan []byte)

	client, err := mqtt.MakeMQTTClient(keypair, endpoint, fmt.Sprintf("tunnel-%s", thingName))
	if err != nil {
		return err
	}
	log.Println("Connected to MQTT")

	// Subscribe to Tunnel Notify topic
	if token := client.Subscribe(
		fmt.Sprintf("$aws/things/%s/tunnels/notify", thingName),
		0,
		func(client paho.Client, msg paho.Message) {
			notifyChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Println("Subscribed to Tunnel Notify topic")

	for {
		select {
		case notifyMsg, ok := <-notifyChan:
			if !ok {
				return fmt.Errorf("failed to read from tunnel notify channel")
			}
			payload := tunnelPayload{}
			if err := json.Unmarshal(notifyMsg, &payload); err != nil {
				return fmt.Errorf("failed to unmarshal tunnel notify message: %w", err)
			}

			if payload.ClientMode != "destination" {
				return fmt.Errorf("tunnel client mode %s is not \"destination\"", payload.ClientMode)
			}

			if len(payload.Services) == 0 {
				return fmt.Errorf("tunnel services are empty")
			}
			if len(payload.Services) > 1 {
				return fmt.Errorf("tunnel services are not a single service")
			}
			if payload.Services[0] != "SSH" {
				return fmt.Errorf("tunnel service %s is not \"SSH\"", payload.Services[0])
			}

			log.Println("TUNNEL REQUESTED")
			return StartLocalProxy(payload.ClientAccessToken, payload.Region)
		}
	}
}
