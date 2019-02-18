package asnlookup

import (
	"errors"
	"strconv"
	"strings"
)

var (
	// ErrInvalidIPv4Address is returned when IPv4 address is badly formatted
	ErrInvalidIPv4Address = errors.New("Invalid IPv4 Address")

	// ErrInvalidIPv4Cidr is returned when IPv4 address & CIDR is badly formatted
	ErrInvalidIPv4Cidr = errors.New("Invalid IPv4 CIDR Format")
)

// IPv4Address struct satisfies IPAddress interface. This is used in trie Insert()
// & Find() functions to make them IP address type agnostic in those functions
type IPv4Address struct {
	cidrLen int
	mask    uint32
	ip      uint32
	ipStr   string
	asn     int
}

// Compile time check to ensure IPv4Address satiesfies IPAddress interface
var _ IPAddress = &IPv4Address{}

// newIPv4Address returns new IPv4Address
func newIPv4Address(ipCidr string, asn int) (IPAddress, error) {
	if isValidIPv4Cidr(ipCidr) == false {
		return nil, ErrInvalidIPv4Cidr
	}

	ipv4Address := IPv4Address{}
	ip := strings.Split(ipCidr, "/")
	ipInt, err := ipv4StrToInt(ip[0])
	if err != nil {
		return ipv4Address, err
	}

	prefix, err := strconv.Atoi(ip[1])
	if err != nil {
		return ipv4Address, err
	}

	ipv4Address.cidrLen = prefix
	ipv4Address.mask = uint32(^(uint32(0))) << uint32(32-prefix)
	ipv4Address.ip = ipInt & ipv4Address.mask
	ipStr, err := intToIPv4Str(ipv4Address.ip)
	if err != nil {
		return ipv4Address, err
	}
	ipv4Address.ipStr = ipStr
	ipv4Address.asn = asn

	return ipv4Address, nil
}

// GetString returns IPv4 address in string format. This method
// is needed to satisfy IPAddress interface
func (ipv4 IPv4Address) GetString() string {
	return ipv4.ipStr
}

// GetNthHighestBit returns nth highest bit for IPv4 address. This method
// is needed to satisfy IPAddress interface
func (ipv4 IPv4Address) GetNthHighestBit(n uint8) uint8 {
	nthBit := ((ipv4.ip) >> (32 - uint32(n))) & 0x1
	return uint8(nthBit)
}

// GetAsn returns ASN stored in IPv4 address. This method
// is needed to satisfy IPAddress interface
func (ipv4 IPv4Address) GetAsn() int {
	return ipv4.asn
}

// GetCidrLen returns CIDR prefix length stored in IPv4 address.
// This method is needed to satisfy IPAddress interface
func (ipv4 IPv4Address) GetCidrLen() int {
	return ipv4.cidrLen
}

// GetNumBitsInAddress returns number of bits in IPv4 address.
// This method is needed to satisfy IPAddress interface
func (ipv4 IPv4Address) GetNumBitsInAddress() int {
	return 32
}

// Following are helper functions to parse, validate & convert IPv4 address
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

func isValidIPv4Cidr(cidr string) bool {
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
func ipv4StrToInt(s string) (uint32, error) {

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

func intToIPv4Str(ip uint32) (string, error) {
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
