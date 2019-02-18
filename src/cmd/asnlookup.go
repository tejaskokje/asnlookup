package main

import (
	"asnlookup"
	"fmt"
	"os"
)

func main() {

	cfg, err := asnlookup.GetConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	nodeInfoList := asnlookup.Find(cfg)
	for _, info := range nodeInfoList {
		fmt.Printf("%s/%d %d\n", info.Subnet, info.Cidr, info.Asn)
	}
}
