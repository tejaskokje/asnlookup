package asnlookup

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	IPCidrList       string
	IPToFind         IPAddress
	NewIPAddressFunc func(string, int) (IPAddress, error)
	IsValidCidrFunc  func(string) bool
}

var (
	ErrNoIPToFind            = errors.New("Please provide IPv4 or IPv6 address for ASN lookup")
	ErrMoreThanOneIPToFind   = errors.New("Please provide only one IP address")
	ErrInvalidInputIPAddress = errors.New("Invalid IP address in input")
)

func GetConfig() (*Config, error) {

	cfg := &Config{}
	args := os.Args[1:]
	numArgs := len(args)
	if numArgs == 0 {
		return nil, ErrNoIPToFind
	} else if numArgs > 1 {
		return nil, ErrMoreThanOneIPToFind
	}

	reqIP := args[0]
	var ipToFind IPAddress
	if IsValidIPv4(reqIP) {
		ipToFind, err := NewIPv4Address(reqIP+"/32", -1)
		if err != nil {
			return nil, err
		}

		cfg.IPToFind = ipToFind
		cfg.NewIPAddressFunc = NewIPv4Address
		cfg.IsValidCidrFunc = IsValidIPv4Cidr
	} else if IsValidIPv6(reqIP) {
		ipToFind, err := NewIPv6Address(reqIP+"/128", -1)
		if err != nil {
			return nil, err
		}
		cfg.IPToFind = ipToFind
		cfg.NewIPAddressFunc = NewIPv6Address
		cfg.IsValidCidrFunc = IsValidIPv6Cidr
	} else {
		return nil, ErrInvalidInputIPAddress
	}

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

	cfg.IPCidrList = string(text)

	return cfg, nil
}
