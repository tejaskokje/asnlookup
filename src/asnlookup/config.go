package asnlookup

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	IPToFind      IPAddress
	IPAddressList []IPAddress
	trie          *Trie
}

var (
	// ErrNoIPToFind is returned when input does not have target IP address to find
	ErrNoIPToFind = errors.New("Please provide IPv4 or IPv6 address for ASN lookup")

	// ErrMoreThanOneIPToFind is returned when there are more than one target IP address in input
	ErrMoreThanOneIPToFind = errors.New("Please provide only one IP address")

	// ErrInvalidInputIPAddress is returned when input IP address (IPv4 or IPv6) is invalid
	ErrInvalidInputIPAddress = errors.New("Invalid IP address in input")
)

// GetConfig generates configuration and creates trie for lookup.
// It uses CONFIG_FILE_PATH environment variable (to get IP, CIDR & ASN information) if defined.
// Otherwise it uses default URL address to fetch configuration from.
// It also gets target IP to lookup from command line arguments.
// It returns a pointer to Config structure which holds all this information.
func GetConfig(envTargetIP ...string) (*Config, error) {

	cfg := &Config{}

	// If target IP is passed in as an argument, then use that.
	// Otherwise get it from environment variable.
	var reqIPStr string
	if len(envTargetIP) > 0 {
		reqIPStr = envTargetIP[0]
	} else {
		args := os.Args[1:]
		numArgs := len(args)
		if numArgs == 0 {
			return nil, ErrNoIPToFind
		} else if numArgs > 1 {
			return nil, ErrMoreThanOneIPToFind
		}
		reqIPStr = args[0]
	}

	var newIPAddressFunc func(string, int) (IPAddress, error)
	var isValidCidrFunc func(string) bool

	// Setup correct information for trie creation based on IP address type
	if isValidIPv4(reqIPStr) {
		ipToFind, err := newIPv4Address(reqIPStr+"/32", -1)
		if err != nil {
			return nil, err
		}

		cfg.IPToFind = ipToFind
		newIPAddressFunc = newIPv4Address
		isValidCidrFunc = isValidIPv4Cidr
	} else if isValidIPv6(reqIPStr) {
		ipToFind, err := newIPv6Address(reqIPStr+"/128", -1)
		if err != nil {
			return nil, err
		}
		cfg.IPToFind = ipToFind
		newIPAddressFunc = newIPv6Address
		isValidCidrFunc = isValidIPv6Cidr
	} else {
		return nil, ErrInvalidInputIPAddress
	}

	cfg.trie = NewTrie()

	// Get configuration either from CONFIG_FILE_PATH environment variable or
	// default configURL
	var reader io.Reader
	configURL := "http://lg01.infra.ring.nlnog.net/table.txt"
	configFile := os.Getenv("CONFIG_FILE_PATH")

	if configFile == "" {
		resp, err := http.Get(configURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		reader = resp.Body
	} else {

		file, err := os.Open(configFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		reader = file
	}

	text, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Scan the text line by line and insert ipAddress information into trie
	scanner := bufio.NewScanner(strings.NewReader(string(text)))
	for scanner.Scan() {
		parts := strings.Split(strings.Trim(scanner.Text(), " "), " ")
		if len(parts) != 2 {
			continue
		}

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

			Insert(cfg.trie, ipAddress)
			cfg.IPAddressList = append(cfg.IPAddressList, ipAddress)
		}
	}

	return cfg, nil
}
