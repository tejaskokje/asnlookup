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
	ErrNoIPToFind            = errors.New("Please provide IPv4 or IPv6 address for ASN lookup")
	ErrMoreThanOneIPToFind   = errors.New("Please provide only one IP address")
	ErrInvalidInputIPAddress = errors.New("Invalid IP address in input")
)

func GetConfig(envTargetIP ...string) (*Config, error) {

	cfg := &Config{}
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

	if IsValidIPv4(reqIPStr) {
		ipToFind, err := NewIPv4Address(reqIPStr+"/32", -1)
		if err != nil {
			return nil, err
		}

		cfg.IPToFind = ipToFind
		newIPAddressFunc = NewIPv4Address
		isValidCidrFunc = IsValidIPv4Cidr
	} else if IsValidIPv6(reqIPStr) {
		ipToFind, err := NewIPv6Address(reqIPStr+"/128", -1)
		if err != nil {
			return nil, err
		}
		cfg.IPToFind = ipToFind
		newIPAddressFunc = NewIPv6Address
		isValidCidrFunc = IsValidIPv6Cidr
	} else {
		return nil, ErrInvalidInputIPAddress
	}

	cfg.trie = NewTrie()

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
