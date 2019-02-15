package asnlookup

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidIPv4Address = errors.New("Invalid IPv4 Address")
)

type IPv4Address [4]byte

func ParseIPv4(ipStr string) (IPv4Address, error) {
	ipv4Addr := IPv4Address{}
	if IsValidIPv4(ipStr) == false {
		return ipv4Addr, ErrInvalidIPv4Address
	}

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

func IPv4ToInt(s string) (uint32, error) {

	var ipv4Int uint32 = 0
	ipAddr, err := ParseIPv4(s)
	if err != nil {
		return ipv4Int, err
	}
	ipv4Int = (ipv4Int | uint32(ipAddr[0])) << 8
	ipv4Int = (ipv4Int | uint32(ipAddr[1])) << 8
	ipv4Int = (ipv4Int | uint32(ipAddr[2])) << 8
	ipv4Int = (ipv4Int | uint32(ipAddr[3]))
	return ipv4Int, nil
}

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

type IPv6Address [16]byte

func ParseIPv6(ipStr string) (IPv6Address, error) {
	ipv6Addr := IPv6Address{}
	/*
		if IsValidIPv4(ipStr) == false {
			return ipv4Addr, ErrInvalidIPv4Address
		}

		parts := strings.Split(ipStr, ".")
		for i, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				return ipv4Addr, err
			}

			ipv4Addr[i] = byte(num)
		}
	*/
	return ipv6Addr, nil
}

func IPv6ToInt(s string) (uint32, error) {

	var ipv6Int uint32 = 0
	/*
		ipAddr, err := ParseIPv4(s)
		if err != nil {
			return ipv4Int, err
		}
		ipv4Int = (ipv4Int | uint32(ipAddr[0])) << 8
		ipv4Int = (ipv4Int | uint32(ipAddr[1])) << 8
		ipv4Int = (ipv4Int | uint32(ipAddr[2])) << 8
		ipv4Int = (ipv4Int | uint32(ipAddr[3]))
	*/
	return ipv6Int, nil
}

func IsValidIPv6(ip string) bool {

	cHextets := strings.Split(ip, "::")
	if len(cHextets) > 2 {
		return false
	}

	if len(cHextets) == 2 {
		leftHextets := strings.Split(cHextets[0], ":")
		rightHextets := strings.Split(cHextets[1], ":")
		if len(leftHextets)+len(rightHextets) > 8 {
			return false
		}

		for _, hextetSlice := range [][]string{leftHextets, rightHextets} {
			for _, hextet := range hextetSlice {
				if len(hextet) > 4 {
					return false
				}

				for _, h := range hextet {
					if (h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') ||
						(h >= 'A' && h <= 'F') {
						continue
					} else {
						return false
					}
				}

			}
		}

	} else {
		uHextets := strings.Split(ip, ":")
		if len(uHextets) != 8 {
			return false
		}

	}

	return true
}

func IsValidIPv6Cidr(cidr string) bool {

	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return false
	}

	if IsValidIPv6(parts[0]) == false {
		return false
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}

	if num < 1 || num > 128 {
		return false
	}

	return true
}
