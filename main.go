package main

import (
	"fmt"
	"log"
	"os"

	"github.com/patrickjmcd/aws-iot-device-sdk-go/cmd"
	"github.com/patrickjmcd/go-version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cmd.RegisterCmd)
	rootCmd.AddCommand(cmd.GetMACAddressCmd)
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
