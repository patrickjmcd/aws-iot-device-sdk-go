// Copyright 2020 SEQSENSE, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/patrickjmcd/aws-iot-device-sdk-go/pkg/tunnel"
	"github.com/patrickjmcd/go-version"
	"github.com/spf13/cobra"
)

var (
	accessToken     string
	proxyEndpoint   string
	region          string
	sourcePort      int
	destinationApp  string
	noSSLHostVerify bool
	proxyScheme     string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&accessToken, "access-token", "", "Client access token")
	rootCmd.PersistentFlags().StringVar(&proxyEndpoint, "proxy-endpoint", "", "Endpoint of proxy server (e.g. data.tunneling.iot.ap-northeast-1.amazonaws.com:443)")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "Endpoint region. Exclusive flag with -proxy-endpoint")
	rootCmd.PersistentFlags().IntVar(&sourcePort, "source-listen-port", 0, "Assigns source mode and sets the port to listen")
	rootCmd.PersistentFlags().StringVar(&destinationApp, "destination-app", "", "Assigns destination mode and set the endpoint in address:port format")
	rootCmd.PersistentFlags().BoolVar(&noSSLHostVerify, "no-ssl-host-verify", false, "Turn off SSL host verification")
	rootCmd.PersistentFlags().StringVar(&proxyScheme, "proxy-scheme", "wss", "Proxy server protocol scheme")
}

var rootCmd = &cobra.Command{
	Use:   "localproxy",
	Short: "run a local proxy server for AWS IoT Core Secure Tunneling",
	Long:  `Run a local proxy server for AWS IoT Core Secure Tunneling`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if accessToken == "" {
			log.Fatal("--access-token is required")
		}

		if proxyEndpoint == "" && region == "" {
			log.Fatal("--proxy-endpoint or --region is required")
		}

		if sourcePort == 0 && destinationApp == "" {
			log.Fatal("--source-listen-port or --destination-app is required")
		}

		params := tunnel.ProxyParams{
			AccessToken:     accessToken,
			ProxyEndpoint:   proxyEndpoint,
			Region:          region,
			SourcePort:      sourcePort,
			DestinationApp:  destinationApp,
			NoSSLHostVerify: noSSLHostVerify,
		}
		err := tunnel.StartLocalProxy(params)
		if err != nil {
			log.Fatal(err)
		}
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

// func main() {
// 	flag.Parse()

// if *accessToken == "" {
// 	log.Fatal("error: -access-token must be specified")
// }

// var endpoint string
// switch {
// case *proxyEndpoint != "" && *region == "":
// 	endpoint = *proxyEndpoint
// case *region != "" && *proxyEndpoint == "":
// 	endpoint = fmt.Sprintf("data.tunneling.iot.%s.amazonaws.com", *region)
// default:
// 	log.Fatal("error: one of -proxy-endpoint or -region must be specified")
// }

// proxyOpts := []tunnel.ProxyOption{
// 	func(opt *tunnel.ProxyOptions) error {
// 		opt.InsecureSkipVerify = *noSSLHostVerify
// 		opt.Scheme = *proxyScheme
// 		return nil
// 	},
// 	tunnel.WithErrorHandler(tunnel.ErrorHandlerFunc(func(err error) {
// 		log.Print(err)
// 	})),
// }

// switch {
// case *sourcePort > 0 && *destinationApp == "":
// 	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *sourcePort))
// 	if err != nil {
// 		log.Fatalf("error: %v", err)
// 	}
// 	err = tunnel.ProxySource(listener, endpoint, *accessToken, proxyOpts...)
// 	if err != nil {
// 		log.Fatalf("error: %v", err)
// 	}

// case *destinationApp != "" && *sourcePort == 0:
// 	err := tunnel.ProxyDestination(func() (io.ReadWriteCloser, error) {
// 		return net.Dial("tcp", *destinationApp)
// 	}, endpoint, *accessToken, proxyOpts...)
// 	if err != nil {
// 		log.Fatalf("error: %v", err)
// 	}

// default:
// 	log.Fatal("error: one of -source-listen-port or -destination-app must be specified")
// }
// }
