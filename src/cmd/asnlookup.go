package main

import (
	"asnlookup"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	cfg, err := asnlookup.GetConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	t := asnlookup.NewTrie()

	scanner := bufio.NewScanner(strings.NewReader(cfg.IPCidrList))
	for scanner.Scan() {
		parts := strings.Split(strings.Trim(scanner.Text(), " "), " ")
		if len(parts) != 2 {
			continue
		}
		isValidCidr := cfg.IsValidCidrFunc(parts[0])
		if isValidCidr == true {
			asn, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}
			ipAddress, err := cfg.NewIPAddressFunc(parts[0], asn)
			if err != nil {
				continue
			}
			asnlookup.Insert(t, ipAddress)
		}
	}

	nodeInfoList := asnlookup.Find(t, ipToFind)
	for _, info := range nodeInfoList {
		fmt.Printf("%s/%d %d\n", info.Subnet, info.Cidr, info.Asn)
	}
}
