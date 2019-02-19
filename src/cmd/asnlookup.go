package main

import (
	"asnlookup"
	"fmt"
	"os"
)

func main() {

	// Get the configuration
	cfg, err := asnlookup.GetConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// Do a lookup
	nodeInfoList := asnlookup.Find(cfg)
	if len(nodeInfoList) == 0 {
		os.Exit(1)
	}

	for _, info := range nodeInfoList {
		fmt.Printf("%s/%d %d\n", info.Subnet, info.Cidr, info.Asn)
	}
}
