package tunnel

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/seqsense/aws-iot-device-sdk-go/v5/tunnel"
)

// ProxyParams holds the parameters for running the local proxy
type ProxyParams struct {
	AccessToken     string
	ProxyEndpoint   string
	Region          string
	SourcePort      int
	DestinationApp  string
	NoSSLHostVerify bool
	ProxyScheme     string
}

// StartLocalProxy starts a local proxy to the AWS IoT MQTT endpoint
func StartLocalProxy(params ProxyParams) error {
	if params.ProxyScheme == "" {
		params.ProxyScheme = "wss"
	}

	if params.AccessToken == "" {
		log.Fatal("error: AccessToken must be specified")
	}

	var endpoint string
	switch {
	case params.ProxyEndpoint != "" && params.Region == "":
		endpoint = params.ProxyEndpoint
	case params.Region != "" && params.ProxyEndpoint == "":
		endpoint = fmt.Sprintf("data.tunneling.iot.%s.amazonaws.com", params.Region)
	default:
		log.Fatal("error: one of ProxyEndpoint or Region must be specified")
	}

	proxyOpts := []tunnel.ProxyOption{
		func(opt *tunnel.ProxyOptions) error {
			opt.InsecureSkipVerify = params.NoSSLHostVerify
			opt.Scheme = params.ProxyScheme
			return nil
		},
		tunnel.WithErrorHandler(tunnel.ErrorHandlerFunc(func(err error) {
			log.Print(err)
		})),
	}

	switch {
	case params.SourcePort > 0 && params.DestinationApp == "":
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", params.SourcePort))
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		err = tunnel.ProxySource(listener, endpoint, params.AccessToken, proxyOpts...)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

	case params.DestinationApp != "" && params.SourcePort == 0:
		err := tunnel.ProxyDestination(func() (io.ReadWriteCloser, error) {
			return net.Dial("tcp", params.DestinationApp)
		}, endpoint, params.AccessToken, proxyOpts...)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

	default:
		log.Fatal("error: one of SourcePort or DestinationApp must be specified")
	}
	return nil
}
