package thing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/google/uuid"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/models"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/mqtt"
	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/networking"
	"github.com/spf13/cobra"
)

var (
	endpoint              string
	templateName          string
	thingName             string
	privateKeyPath        string
	certificatePath       string
	rootCAPath            string
	outputFilePath        string
	parameterJSONFilePath string
	clientID              string
	parameters            map[string]string
)

func init() {
	RegisterCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "The AWS IoT endpoint")
	RegisterCmd.PersistentFlags().StringVarP(&templateName, "template", "t", "", "The template to use")
	RegisterCmd.PersistentFlags().StringVarP(&privateKeyPath, "private-key", "k", "", "The private key path")
	RegisterCmd.PersistentFlags().StringVarP(&certificatePath, "certificate", "c", "", "The certificate path")
	RegisterCmd.PersistentFlags().StringVarP(&rootCAPath, "root-ca", "r", "", "The root CA path")
	RegisterCmd.PersistentFlags().StringVarP(&outputFilePath, "output", "o", ".", "The output file path")
	RegisterCmd.PersistentFlags().StringVarP(&parameterJSONFilePath, "parameters", "p", "", "The parameters file path")
	RegisterCmd.PersistentFlags().StringVarP(&clientID, "client-id", "i", "", "The client ID")
}

func checkRegisterParameters() error {
	if endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	if templateName == "" {
		return fmt.Errorf("template is required")
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
	if parameterJSONFilePath != "" {
		file, err := ioutil.ReadFile(parameterJSONFilePath)
		if err != nil {
			log.Fatalf("error reading parameters file: %v", err)
		}

		err = json.Unmarshal(file, &parameters)
		if err != nil {
			log.Fatalf("error unmarshalling parameters file: %v", err)
		}
	}

	if clientID == "" {
		clientID = uuid.New().String()
	}

	return nil
}

// RegisterCmd registers a new thing
var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers a new thing",
	Long:  `Registers a new thing and creates/stores all the keys and certs needed to communicate with AWS IoT Core`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkRegisterParameters(); err != nil {
			log.Fatal(err)
		}

		macAddress, _, err := networking.GetMACAddress()
		if err != nil {
			log.Fatal(err)
		}
		uniqueID := macAddress[:6] + "fffe" + macAddress[6:]
		parameters["UniqueId"] = string(uniqueID)

		keypair := models.KeyPair{
			PrivateKeyPath:    privateKeyPath,
			CertificatePath:   certificatePath,
			CACertificatePath: rootCAPath,
		}

		client, err := mqtt.MakeMQTTClient(keypair, endpoint, clientID)
		if err != nil {
			log.Fatal(err)
		}

		err = ProvisionThing(client, keypair, endpoint, templateName, parameters, outputFilePath)
		if err != nil {
			log.Fatal(err)
		}
	},
}
