package utils

import (
	"fmt"
	"log"
	"net"
)

// GetLocalIp obtains server's local IP address by using the routing mechanism's determination
// of what IP address it would use to reach remote IP (8.8.8.8:80 = Google's public DNS server).
func GetLocalIp() (string, error) {
	connection, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return "", err
	}

	localAddr, ok := connection.LocalAddr().(*net.UDPAddr)

	if !ok {
		return "", fmt.Errorf("type assert failed")
	}

	if err := connection.Close(); err != nil {
		log.Printf("Error while closing connection for GetLocalIp(): %s", err)
	}

	return localAddr.IP.String(), nil
}
