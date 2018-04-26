package internals

import (
	"fmt"
	"net"
	"os"
)

func GetIPAddress() string {

	netInterfaceAddresses, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, netInterfaceAddress := range netInterfaceAddresses {

		networkIP, ok := netInterfaceAddress.(*net.IPNet)

		if ok && !networkIP.IP.IsLoopback() && networkIP.IP.To4() != nil {

			ip := networkIP.IP.String()
			return ip
		}
	}
	return ""
}

func GetRPCPort() string {
	p, ok := os.LookupEnv("RPC_PORT")
	if ok {
		return fmt.Sprintf(":%s", p)
	}
	return ":5001"
}

func GetEndpoint() string {
	return GetIPAddress() + GetRPCPort()
}
