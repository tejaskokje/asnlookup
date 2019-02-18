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

	var newIPAddressFunc func(string, int) (asnlookup.IPAddress, error)
	var isValidCidrFunc func(string) bool
	var ipToFind asnlookup.IPAddress
	if asnlookup.IsValidIPv4(cfg.IPToFind) {
		//Check if reqIP is valid IPV4 or v6
		ipToFind, err = asnlookup.NewIPv4Address(cfg.IPToFind+"/32", -1)
		if err != nil {
			fmt.Printf("Error encountered for target IP %s: %v\n", cfg.IPToFind, err)
			os.Exit(1)
		}
		newIPAddressFunc = asnlookup.NewIPv4Address
		isValidCidrFunc = asnlookup.IsValidIPv4Cidr
	} else if asnlookup.IsValidIPv6(cfg.IPToFind) {
		ipToFind, err = asnlookup.NewIPv6Address(cfg.IPToFind+"/128", -1)
		if err != nil {
			fmt.Printf("Error encountered for target IP %s: %v\n", cfg.IPToFind, err)
			os.Exit(1)
		}
		newIPAddressFunc = asnlookup.NewIPv6Address
		isValidCidrFunc = asnlookup.IsValidIPv6Cidr

	} else {
		fmt.Println("Invalid IP address in input")
		os.Exit(1)
	}

	t := asnlookup.NewTrie()

	scanner := bufio.NewScanner(strings.NewReader(cfg.IPCidrList))
	for scanner.Scan() {
		parts := strings.Split(strings.Trim(scanner.Text(), " "), " ")
		isValidCidr := isValidCidrFunc(parts[0])
		if isValidCidr == true {
			asn, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}
			ipAddress, err := newIPAddressFunc(parts[0], asn)
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
