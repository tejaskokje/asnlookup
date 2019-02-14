package main

import (
	"asnlookup"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	configURL = "http://lg01.infra.ring.nlnog.net/table.txt"
)

func getConfig() (string, error) {
	var reader io.Reader
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

	t := asnlookup.NewTrie()

	scanner := bufio.NewScanner(strings.NewReader(cfgStr))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "#") {
			continue
		}

		parts := strings.Split(strings.Trim(scanner.Text(), " "), " ")
		isValidCidr := asnlookup.IsValidIPv4Cidr(parts[0])
		if isValidCidr == false {
			//fmt.Println("Invalid input: ", parts[0])
		} else {
			ip := strings.Split(parts[0], "/")
			//fmt.Println(asnlookup.IPv4ToInt(net.ParseIP(ip[0])))
			ipInt := asnlookup.IPv4ToInt(net.ParseIP(ip[0]))
			prefix, _ := strconv.Atoi(ip[1])
			mask := uint32(^(uint32(0))) << uint32(32-prefix)
			asn, _ := strconv.Atoi(parts[1])
			//fmt.Printf("%032b\n", ipInt)
			//fmt.Printf("%032b\n", mask)
			//fmt.Println(ipInt, ipInt&mask, prefix, asn)
			asnlookup.Insert(t, ipInt&mask, asn, prefix)
		}
	}

	//asnlookup.Insert(t, 12, 50, 4)
	fmt.Println(asnlookup.Find(t, 134744072))
	//fmt.Println("Dump Trie")
	//asnlookup.DumpTrie(t)

}
