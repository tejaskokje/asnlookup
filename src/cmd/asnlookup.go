package main

import (
	"asnlookup"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getConfig() (string, error) {
	var reader io.Reader

	configURL := "http://lg01.infra.ring.nlnog.net/table.txt"
	configFile := os.Getenv("CONFIG_FILE_PATH")

	if configFile == "" {
		resp, err := http.Get(configURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		reader = resp.Body
	} else {

		file, err := os.Open(configFile)
		if err != nil {
			return "", err
		}
		defer file.Close()

		reader = file
	}

	text, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(text), nil

}

func main() {

	cfgStr, err := getConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	args := os.Args[1:]
	numArgs := len(args)
	if numArgs == 0 {
		fmt.Println("Please provide IPv4 or IPv6 address for ASN lookup")
		os.Exit(1)
	} else if numArgs > 1 {
		fmt.Println("Please provide only one IP address")
		os.Exit(1)
	}

	reqIP := args[0]
	var newIPAddressFunc func(string, int) (asnlookup.IPAddress, error)
	var isValidCidrFunc func(string) bool
	var ipToFind asnlookup.IPAddress
	if asnlookup.IsValidIPv4(reqIP) {
		//Check if reqIP is valid IPV4 or v6
		ipToFind, err = asnlookup.NewIPv4Address(reqIP+"/32", -1)
		if err != nil {
			fmt.Printf("Error encountered for target IP %s: %v\n", reqIP, err)
			os.Exit(1)
		}
		newIPAddressFunc = asnlookup.NewIPv4Address
		isValidCidrFunc = asnlookup.IsValidIPv4Cidr
	} else if asnlookup.IsValidIPv6(reqIP) {
		ipToFind, err = asnlookup.NewIPv6Address(reqIP+"/128", -1)
		if err != nil {
			fmt.Printf("Error encountered for target IP %s: %v\n", reqIP, err)
			os.Exit(1)
		}
		newIPAddressFunc = asnlookup.NewIPv6Address
		isValidCidrFunc = asnlookup.IsValidIPv6Cidr

	} else {
		fmt.Println("Invalid IP address in input")
		os.Exit(0)
	}

	t := asnlookup.NewTrie()

	scanner := bufio.NewScanner(strings.NewReader(cfgStr))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "#") {
			continue
		}

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

	//asnlookup.DumpTrie(t)

}
