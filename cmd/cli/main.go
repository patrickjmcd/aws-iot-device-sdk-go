package provisioning

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/patrickjmcd/go-version"
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

var rootCmd = &cobra.Command{
	Use:   "aws-provision",
	Short: "AWS Provisioning Tools",
	Long:  `Complete documentation is available at http://github.com/meshifyiot/aws-iot-core-register-api`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("you didn't specify a command")
	},
}

// execute adds all child commands to the root command sets flags appropriately.
func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	log.Printf(
		"Starting the service...\ncommit: %s, build time: %s, release: %s",
		version.Commit, version.BuildTime, version.Release,
	)
	execute()
}
