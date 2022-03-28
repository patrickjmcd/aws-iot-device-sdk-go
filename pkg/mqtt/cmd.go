package mqtt

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

var (
	endpoint        string
	privateKeyPath  string
	certificatePath string
	rootCAPath      string
	clientID        string
)

func init() {

	CheckCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "The AWS IoT endpoint")
	CheckCmd.PersistentFlags().StringVarP(&privateKeyPath, "private-key", "k", "", "The private key path")
	CheckCmd.PersistentFlags().StringVarP(&certificatePath, "certificate", "c", "", "The certificate path")
	CheckCmd.PersistentFlags().StringVarP(&rootCAPath, "root-ca", "r", "", "The root CA path")
	CheckCmd.PersistentFlags().StringVarP(&clientID, "client-id", "i", "", "The client ID to use")
}

func checkParameters() error {
	if endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}

	if privateKeyPath == "" {
		return fmt.Errorf("private-key is required")
	}
	if certificatePath == "" {
		return fmt.Errorf("certificate is required")
	}
	if rootCAPath == "" {
		return fmt.Errorf("root-ca is required")
	}
	if clientID == "" {
		clientID = uuid.New().String()
	}
	return nil
}

// CheckCmd checks if an MQTT connnection can be established
var CheckCmd = &cobra.Command{
	Use:   "check-mqtt",
	Short: "Checks if an MQTT connection can be established",
	Long:  `Checks if an MQTT connectoin can be established with the provided certs`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkParameters(); err != nil {
			log.Fatal(err)
		}

		keypair := models.KeyPair{
			PrivateKeyPath:    privateKeyPath,
			CertificatePath:   certificatePath,
			CACertificatePath: rootCAPath,
		}

		if _, err := MakeMQTTClient(keypair, endpoint, clientID); err != nil {
			log.Fatal(err)
		}

		log.Printf("MQTT connection established")
	},
}
