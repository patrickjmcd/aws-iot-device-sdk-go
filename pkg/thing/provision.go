package thing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/models"
)

// CreateKeysAndCertificateAcceptedCh holds the bytes of a CreateKeysAndCertificateAccepted message.
type CreateKeysAndCertificateAcceptedCh []byte

// CreateKeysAndCertificateAccepted holds the data from an accepted request to create a new key and certificate.
type CreateKeysAndCertificateAccepted struct {
	CertificateID             string `json:"certificateId"`
	CertificatePem            string `json:"certificatePem"`
	PrivateKey                string `json:"privateKey"`
	CertificateOwnershipToken string `json:"certificateOwnershipToken"`
}

// AWSMQTTErrorCh holds the bytes of a CreateKeysAndCertificateRejected message.
type AWSMQTTErrorCh []byte

// AWSMQTTError holds the data from a rejected request
type AWSMQTTError struct {
	StatusCode   int    `json:"statusCode"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// RegisterThingRequest holds the values needed to make a request to register a new thing.
type RegisterThingRequest struct {
	TemplateName              string            `json:"-"`
	CertificateOwnershipToken string            `json:"certificateOwnershipToken"`
	Parameters                map[string]string `json:"parameters"`
}

// RegisterThingResponse holds the values returned from a successful registration request.
type RegisterThingResponse struct {
	ThingName           string            `json:"thingName"`
	DeviceConfiguration map[string]string `json:"deviceConfiguration"`
}

// RegisterThingAcceptedCh holds the bytes of a RegisterThingAccepted message.
type RegisterThingAcceptedCh []byte

func registerThing(c mqtt.Client, templateName string, certificateOwnershipToken string, parameters map[string]string) error {
	req := RegisterThingRequest{
		TemplateName:              templateName,
		CertificateOwnershipToken: certificateOwnershipToken,
		Parameters:                parameters,
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal register thing request: %v", err)
	}

	if token := c.Publish(
		fmt.Sprintf("$aws/provisioning-templates/%s/provision/json", templateName),
		0,
		false,
		[]byte(reqJSON),
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func writeCertificateFiles(certs CreateKeysAndCertificateAccepted, outputFilePath string) error {

	if err := ioutil.WriteFile(outputFilePath+"/device.certificate.pem", []byte(certs.CertificatePem), 0644); err != nil {
		return fmt.Errorf("failed to write device.cert.pem: %v", err)
	}

	if err := ioutil.WriteFile(outputFilePath+"/device.private.key", []byte(certs.PrivateKey), 0644); err != nil {
		return fmt.Errorf("failed to write device.private.key: %v", err)
	}

	certJSON, err := json.MarshalIndent(certs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal certs: %v", err)
	}
	if err := ioutil.WriteFile(outputFilePath+"/cert.json", certJSON, 0644); err != nil {
		return fmt.Errorf("failed to write cert.json: %v", err)
	}

	return nil
}

// ProvisionThing creates a new set of certificates for the device
func ProvisionThing(c mqtt.Client, keyPair models.KeyPair, awsEndpoint, templateName string, thingParameters map[string]string, certificateOutputPath string) error {
	// make some channels for the data we need
	createErrorChan := make(chan AWSMQTTErrorCh)
	createAcceptedChan := make(chan CreateKeysAndCertificateAcceptedCh)
	registerErrorChan := make(chan AWSMQTTErrorCh)
	registerAcceptedChan := make(chan RegisterThingAcceptedCh)

	// Subscribe to CreateKeysAndCertificate Accepted topic
	if token := c.Subscribe(
		fmt.Sprintf("$aws/certificates/create/json/accepted"),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			createAcceptedChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	// Subscribe to CreateKeysAndCertificate Rejected topic
	if token := c.Subscribe(
		fmt.Sprintf("$aws/certificates/create/json/rejected"),
		0,
		func(client mqtt.Client, msg mqtt.Message) {
			createErrorChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	// Subscribe to RegisterThing Accepted topic
	if token := c.Subscribe(
		fmt.Sprintf("$aws/provisioning-templates/%s/provision/json/accepted", templateName),
		0, func(client mqtt.Client, msg mqtt.Message) {
			registerAcceptedChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	// Subscribe to RegisterThing Rejected topic
	if token := c.Subscribe(
		fmt.Sprintf("$aws/provisioning-templates/%s/provision/json/rejected", templateName),
		0, func(client mqtt.Client, msg mqtt.Message) {
			registerErrorChan <- msg.Payload()
		},
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	// Publish to CreateKeysAndCertificate topic
	if token := c.Publish(
		fmt.Sprintf("$aws/certificates/create/json"),
		0,
		false,
		[]byte("{}"),
	); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for {
		select {
		case accepted, ok := <-createAcceptedChan:
			if !ok {
				return errors.New("failed to read from create accepted channel")
			}
			createAccepted := CreateKeysAndCertificateAccepted{}
			if err := json.Unmarshal(accepted, &createAccepted); err != nil {
				return fmt.Errorf("failed to unmarshal create accepted: %w", err)
			}
			if createAccepted.CertificateOwnershipToken == "" {
				return errors.New("certificate ownership token is empty")
			}
			err := writeCertificateFiles(createAccepted, certificateOutputPath)
			if err != nil {
				return fmt.Errorf("failed to write certificate files: %w", err)
			}
			err = registerThing(c, templateName, createAccepted.CertificateOwnershipToken, thingParameters)
			if err != nil {
				return fmt.Errorf("failed to register thing: %w", err)
			}
		case err, ok := <-createErrorChan:
			if !ok {
				return errors.New("failed to read from create error channel")
			}
			createError := AWSMQTTError{}
			if err := json.Unmarshal(err, &createError); err != nil {
				return fmt.Errorf("failed to unmarshal create error: %w", err)
			}
			return fmt.Errorf(createError.ErrorMessage)
		case accepted, ok := <-registerAcceptedChan:
			if !ok {
				return errors.New("failed to read from register accepted channel")
			}
			registerAccepted := RegisterThingResponse{}
			if err := json.Unmarshal(accepted, &registerAccepted); err != nil {
				return fmt.Errorf("failed to unmarshal register accepted: %w", err)
			}
			log.Printf("Registered thing: %s\n", registerAccepted.ThingName)
			return nil
		case err, ok := <-registerErrorChan:
			if !ok {
				return errors.New("failed to read from register error channel")
			}
			registerError := AWSMQTTError{}
			if err := json.Unmarshal(err, &registerError); err != nil {
				return fmt.Errorf("failed to unmarshal create error: %w", err)
			}
			return fmt.Errorf(registerError.ErrorMessage)
		}
	}
}
