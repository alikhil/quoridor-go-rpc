package internals

import (
	"fmt"
	"log"
	"net"
	"os"
)

func GetIPAddress() string {

	conn, err := net.Dial("udp", "8.8.8.8:80")

	defer conn.Close()
	if err != nil {
		log.Printf("UTILS: failed to get ip addrres")
		return ""
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
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
