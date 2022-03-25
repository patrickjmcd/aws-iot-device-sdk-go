package provisioning

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

// CreateCertsAndKeysCmd is the command to create certs and keys
var CreateCertsAndKeysCmd = &cobra.Command{
	Use:   "create-cert-and-keys",
	Short: "Create Certificate and Keys needed for device provisioning",
	Long:  `Create the certificate and keys needed for just-in-time device provisioning`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputPath := args[0]

		ctx := context.Background()

		err := prepProvisioningCertificateAndKey(ctx, outputPath)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Certificate and Keys created in %s", outputPath)
	},
}
