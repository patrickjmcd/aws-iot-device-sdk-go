package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/patrickjmcd/go-version"
	"github.com/spf13/cobra"
)

var (
	thingName string
)

func init() {
	openTunnelCmd.Flags().StringVarP(&thingName, "thing-name", "n", "", "thing name")

	rootCmd.AddCommand(createCertsAndKeysCmd)
	rootCmd.AddCommand(openTunnelCmd)

}

// createCertsAndKeysCmd is the command to create certs and keys
var createCertsAndKeysCmd = &cobra.Command{
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

var openTunnelCmd = &cobra.Command{
	Use:   "open-tunnel",
	Short: "Open a tunnel to the AWS IoT Core",
	Long:  `Open a tunnel to the AWS IoT Core`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(thingName) == 0 {
			log.Fatal("--thing-name is required")
		}

		err := createTunnel(thingName)
		if err != nil {
			log.Fatal(err)
		}
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
