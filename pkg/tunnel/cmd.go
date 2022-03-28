package tunnel

import (
	"fmt"
	"log"

	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

var (
	endpoint        string
	privateKeyPath  string
	certificatePath string
	rootCAPath      string
	thingName       string
)

func init() {
	ListenForTunnelCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "The AWS IoT endpoint")
	ListenForTunnelCmd.PersistentFlags().StringVarP(&thingName, "thingname", "t", "", "The thing name to use")
	ListenForTunnelCmd.PersistentFlags().StringVarP(&privateKeyPath, "private-key", "k", "", "The private key path")
	ListenForTunnelCmd.PersistentFlags().StringVarP(&certificatePath, "certificate", "c", "", "The certificate path")
	ListenForTunnelCmd.PersistentFlags().StringVarP(&rootCAPath, "root-ca", "r", "", "The root CA path")

}

func checkParameters() error {
	if endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	if thingName == "" {
		return fmt.Errorf("thingname is required")
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

	return nil
}

// ListenForTunnelCmd listens for a new tunnel to be requested
var ListenForTunnelCmd = &cobra.Command{
	Use:   "listen-for-tunnel",
	Short: "Listens for a tunnel creation notification",
	Long:  `Listens for a tunnel creation notification`,
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

		if err := ListenForTunnel(thingName, keypair, endpoint); err != nil {
			log.Fatal(err)
		}
	},
}
