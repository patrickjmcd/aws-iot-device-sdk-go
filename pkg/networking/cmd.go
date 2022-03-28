package networking

import (
	"log"

	"github.com/spf13/cobra"
)

// GetMACAddressCmd gets the MAC address of the device
var GetMACAddressCmd = &cobra.Command{
	Use:   "get-mac-address",
	Short: "Gets the MAC address of the device",
	Long:  `Gets the MAC address of the device`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		macAddress, interfaceName, err := GetMACAddress()
		if err != nil {
			log.Fatal(err)
		}
		uniqueID := macAddress[:6] + "fffe" + macAddress[6:]
		log.Printf("MAC Address: %s\n", macAddress)
		log.Printf("Unique ID: %s\n", uniqueID)
		log.Printf("Interface Name: %s\n", interfaceName)
	},
}
