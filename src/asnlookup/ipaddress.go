package asnlookup

import (
	"math/big"
	"net"
	"strconv"
	"strings"
)

func IsValidIPv4(ip string) bool {
	octets := strings.Split(ip, ".")

	if len(octets) < 3 {
		return false
	}

	for _, o := range octets {
		num, err := strconv.Atoi(o)
		if err != nil {
			return false
		}

		if num < 0 || num > 255 {
			return false
		}
	}

	return true
}

func IsValidIPv4Cidr(cidr string) bool {
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return false
	}

	if IsValidIPv4(parts[0]) == false {
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

func IPv4ToInt(IPv4Address net.IP) uint32 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address.To4())
	return uint32(IPv4Int.Int64())
}
