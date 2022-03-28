package tunnel

import (
	"log"
	"os/exec"
)

// StartLocalProxy starts a local proxy to the AWS IoT MQTT endpoint
func StartLocalProxy(accessToken, region string) error {
	log.Println("localproxy",
		"-access-token", accessToken,
		"-region", region,
		"-destination-app", "localhost:22")
	cmnd := exec.Command("localproxy",
		"-access-token", accessToken,
		"-region", region,
		"-destination-app", "localhost:22")
	err := cmnd.Run()
	return err
}
