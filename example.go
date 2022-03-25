package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/device"
)

func example() {
	endpoint := "AWS ENDPOINT"
	templateName := "TEMPLATE NAME"
	privateKeyPath := "./private.key"
	certificatePath := "./certificate.pem"
	rootCAPath := "./root-ca.pem"
	outputFilePath := "./output"
	parameters := map[string]string{
		"UniqueID": "12345",
	}

	keypair := device.KeyPair{
		PrivateKeyPath:    privateKeyPath,
		CertificatePath:   certificatePath,
		CACertificatePath: rootCAPath,
	}

	clientID := uuid.New().String()
	client, err := device.MakeMQTTClient(keypair, endpoint, clientID)
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}

	err = device.ProvisionThing(*client, keypair, endpoint, templateName, parameters, outputFilePath)
	if err != nil {
		log.Fatal(err)
	}

}
