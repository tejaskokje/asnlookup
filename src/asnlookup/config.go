package asnlookup

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	IPCidrList string
	IPToFind   string
}

var (
	ErrNoIPToFind          = errors.New("Please provide IPv4 or IPv6 address for ASN lookup")
	ErrMoreThanOneIPToFind = errors.New("Please provide only one IP address")
)

func GetConfig() (*Config, error) {
	args := os.Args[1:]
	numArgs := len(args)
	if numArgs == 0 {
		return nil, ErrNoIPToFind
	} else if numArgs > 1 {
		return nil, ErrMoreThanOneIPToFind
	}

	reqIP := args[0]

	var reader io.Reader
	cfg := &Config{}
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
	cfg.IPToFind = reqIP

	return cfg, nil
}
