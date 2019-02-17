package asnlookup

import (
	"errors"
	"strconv"
	"strings"
)

type IPAddress interface {
	IsValidIPCidr(cidr string)
}

var (
	ErrInvalidIPv4Address = errors.New("Invalid IPv4 Address")
)

func parseIPv4(ipStr string) ([]byte, error) {

	if isValidIPv4(ipStr) == false {
		return []byte{}, ErrInvalidIPv4Address
	}

	ipv4Addr := make([]byte, 4)
	parts := strings.Split(ipStr, ".")
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return ipv4Addr, err
		}

		ipv4Addr[i] = byte(num)
	}

	return ipv4Addr, nil
}

func isValidIPv4(ip string) bool {
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return false
	}

	for idx, o := range octets {
		num, err := strconv.Atoi(o)
		if err != nil {
			return false
		}

		if num < 0 || num > 255 {
			return false
		}

		if idx == 0 && num == 0 {
			return false
		}
	}

	return true
}

func IPv4ToInt(s string) (uint32, error) {

	var ipv4Int uint32 = 0
	ipAddr, err := parseIPv4(s)
	if err != nil {
		return ipv4Int, err
	}
	ipv4Int = (ipv4Int | uint32(ipAddr[0])) << 8
	ipv4Int = (ipv4Int | uint32(ipAddr[1])) << 8
	ipv4Int = (ipv4Int | uint32(ipAddr[2])) << 8
	ipv4Int = (ipv4Int | uint32(ipAddr[3]))
	return ipv4Int, nil
}

func IntToIPv4(ip uint32) (string, error) {
	var ipStr []string
	for i := 1; i <= 4; i++ {
		octet := (ip >> uint32(32-8*i)) & 0xFF
		if (i == 1 && (octet < 1 || octet > 255)) ||
			(octet < 0 || octet > 255) {
			return "", ErrInvalidIPv4Address
		}

		s := strconv.Itoa(int(octet))
		ipStr = append(ipStr, s)
	}

	return strings.Join(ipStr, "."), nil

}

func IsValidIPv4Cidr(cidr string) bool {
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return false
	}

	if isValidIPv4(parts[0]) == false {
		return false
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}

	if num < 1 || num > 32 {
		return false
	}

	return true
}
