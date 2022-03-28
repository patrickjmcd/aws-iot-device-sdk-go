package networking

import (
	"fmt"
	"net"
	"strings"
)

type macAddress string
type interfaceName string

// GetMACAddress returns the MAC address of the current machine.
func GetMACAddress() (macAddress, interfaceName, error) {

	//----------------------
	// Get the local machine IP address
	// https://www.socketloop.com/tutorials/golang-how-do-I-get-the-local-ip-non-loopback-address
	//----------------------

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", "", fmt.Errorf("error getting local IP address: %v", err)
	}

	var currentIP, currentNetworkHardwareName, eth0IP string

	ipv4Addresses := []string{}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				currentIP = ipnet.IP.String()
				ipv4Addresses = append(ipv4Addresses, currentIP)
			}
		}
	}

	interfacesWithIPv4Addresses := []string{}
	// get all the system's or local machine's network interfaces
	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {
		for _, ip := range ipv4Addresses {
			if addrs, err := interf.Addrs(); err == nil {
				for _, addr := range addrs {
					// only interested in the name with current IP address
					if strings.Contains(addr.String(), ip) {
						currentNetworkHardwareName = interf.Name
						interfacesWithIPv4Addresses = append(interfacesWithIPv4Addresses, currentNetworkHardwareName)
						if currentNetworkHardwareName == "eth0" {
							eth0IP = ip
						}
					}
				}
			}
		}
	}

	if len(interfacesWithIPv4Addresses) == 0 {
		return "", "", fmt.Errorf("no ipv4 network interface found")
	}

	if eth0IP != "" {
		// default to using eth0 if it's present
		currentNetworkHardwareName = "eth0"
	}

	// extract the hardware information base on the interface name
	// capture above
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)
	if err != nil {
		return "", "", fmt.Errorf("error getting interface by name: %v", err)
	}

	name := netInterface.Name
	niHardwareAddr := netInterface.HardwareAddr

	// verify if the MAC address can be parsed properly
	hwAddr, err := net.ParseMAC(niHardwareAddr.String())
	if err != nil {
		return "", "", fmt.Errorf("error parsing MAC address: %v", err)
	}
	macString := strings.Replace(hwAddr.String(), ":", "", -1)

	return macAddress(macString), interfaceName(name), nil

}
