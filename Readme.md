# AWS IoT SDK for Go lang

The aws-iot-device-sdk-go package allows developers to write Go lang applications which access the AWS IoT Platform via MQTT.

## Install

`go get "github.com/patrickjmcd/aws-iot-device-sdk-go"`

## Example

```go
package main

import (
    "log"

    "github.com/google/uuid"
    "github.com/patrickjmcd/aws-iot-device-sdk-go/device"
)

func main() {
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
```
