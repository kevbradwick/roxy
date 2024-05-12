package roxy

import (
	"log"
	"net"
)

func IPAllowed(address string, allowList []string) bool {
	ipToCheck := net.ParseIP(address)
	if ipToCheck == nil {
		log.Printf("IP '%s' is not a valid IP4 or IP6 address", ipToCheck)
		return false
	}

	for _, v := range allowList {
		// this is a simple IP check
		if entry := net.ParseIP(v); entry != nil {
			log.Printf("IP '%s' is allowed to pass through based on rule matching '%s'", ipToCheck, entry)
			return ipToCheck.Equal(entry)
		}

		// check CIDR range
		_, ipnet, err := net.ParseCIDR(v)
		if err != nil {
			log.Printf("Failed to parse CIDR '%s', err=%v", v, err)
			return false
		}

		return ipnet.Contains(ipToCheck)
	}
	return false
}
