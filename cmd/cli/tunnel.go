package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iotsecuretunneling"
	"github.com/aws/aws-sdk-go-v2/service/iotsecuretunneling/types"
)

func createTunnel(thingName string) error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := iotsecuretunneling.NewFromConfig(cfg)

	// create a new tunnel
	openOutput, err := client.OpenTunnel(ctx, &iotsecuretunneling.OpenTunnelInput{
		Description: aws.String("aws iot secure tunnel"),
		DestinationConfig: &types.DestinationConfig{
			Services:  []string{"SSH"},
			ThingName: &thingName,
		},
	})
	if err != nil {
		return err
	}

	return localproxyTunnel(cfg.Region, *openOutput.SourceAccessToken)
}

func localproxyTunnel(region, sourceToken string) error {

	_, err := exec.LookPath("localproxy")
	if err != nil {
		return fmt.Errorf("localproxy not found in path")
	}

	log.Println("localproxy",
		"-access-token", sourceToken,
		"-region", region,
		"-source-listen-port", "2222")
	cmnd := exec.Command("localproxy",
		"-access-token", sourceToken,
		"-region", region,
		"-source-listen-port", "2222")
	err = cmnd.Run()
	return err
}
