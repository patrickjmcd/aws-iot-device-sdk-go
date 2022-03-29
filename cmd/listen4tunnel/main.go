package main

import (
	"fmt"
	"log"
	"os"

	"github.com/patrickjmcd/aws-iot-device-sdk-go/cmd/listen4tunnel/cfg"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/models"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/networking"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/tunnel"
)

func main() {

	if cfg.Endpoint == "" {
		if ep := os.Getenv("AWS_IOT_ENDPOINT"); ep != "" {
			cfg.Endpoint = ep
		} else {
			log.Fatalf("endpoint is required")
		}
	}

	keypair := models.KeyPair{
		PrivateKeyPath:    cfg.PrivateKeyPath,
		CertificatePath:   cfg.CertificatePath,
		CACertificatePath: cfg.RootCAPath,
	}

	if cfg.ThingName == "" {
		uniqueID, _, err := networking.GetMACAddress()
		if err != nil {
			log.Fatalf("error getting MAC address: %v", err)
		}
		cfg.ThingName = fmt.Sprintf("gateway-%sfffe%s", uniqueID[0:6], uniqueID[6:])
	}

	log.Println("Endpoint:", cfg.Endpoint)
	log.Println("PrivateKeyPath:", keypair.PrivateKeyPath)
	log.Println("CertificatePath:", keypair.CertificatePath)
	log.Println("CACertificatePath:", keypair.CACertificatePath)
	log.Println("ThingName:", cfg.ThingName)

	err := tunnel.ListenForTunnel(cfg.ThingName, keypair, cfg.Endpoint)
	if err != nil {
		log.Fatalf("error listening for tunnel: %v", err)
	}
	log.Println("SHUTDOWN WITHOUT ERROR")
}
