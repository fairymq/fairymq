package main

import (
	"flag"
	"fmt"
	"net"
)

type Config struct {
	BindAddress          string
	BindPort             uint
	MemberlistPort       uint
	GenerateQueueKeyPair string
	JoinAddresses        []string
}

func GetConfig() Config {
	var joinAddresses []string
	flag.Func(
		"join-address",
		"IP address and memberlist port of a peer in a cluster we would like to join. This flag can be specified multiple times.   --join-address=<ipaddress>:<memberlist-port>",
		func(address string) error {
			if net.ParseIP(address) == nil {
				return fmt.Errorf("invalid ip address: %s", address)
			}
			joinAddresses = append(joinAddresses, address)
			return nil
		})

	bindAddress := flag.String("bind-address", "0.0.0.0", "The host address to bind to.   --bind-address=0.0.0.0")
	bindPort := flag.Uint("bind-port", 5991, "The port to bind to.   --bind-port=5991")
	memberlistPort := flag.Uint("memberlist-port", 7946, "Port used by by this node to communicate with other nodes in the cluster.   --memberlist-port=7946")
	generateQueueKeyPair := flag.String("generate-queue-key-pair", "", "Generates a new queue keypair.   --generate-queue-key-pair=yourqueuename")
	flag.Parse()

	config := Config{
		BindAddress:          *bindAddress,
		BindPort:             *bindPort,
		MemberlistPort:       *memberlistPort,
		GenerateQueueKeyPair: *generateQueueKeyPair,
		JoinAddresses:        joinAddresses,
	}

	return config
}
