package util

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

// GetLocalIP returns the local ip4 address
func GetLocalIP() (net.IP, error) {
	inter, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}
	for _, iface := range inter {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			adrrVal, ok := addr.(*net.IPNet)
			if ok {
				ip = adrrVal.IP
				if ip.To4() != nil {
					if strings.HasPrefix(ip.To4().String(), "192.168.1") {
						return ip, nil
					}
				}
			}

		}

	}
	return nil, errors.New("couldnt find the right ethernet interface.")
}
