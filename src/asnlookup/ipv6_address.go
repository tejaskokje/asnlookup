package asnlookup

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidIPv6Address = errors.New("Invalid IPv6 Address")
)

type IPv6Address []byte

func ParseIPv6(ipStr string) (IPv6Address, error) {
	//var bytes []byte
	bytes := make([]byte, 39)

	if IsValidIPv6(ipStr) == false {
		return IPv6Address{}, ErrInvalidIPv6Address
	}

	ipStr = strings.ToLower(ipStr)
	if strings.Contains(ipStr, "::") {
		parts := strings.Split(ipStr, "::")

		byteCount := 0
		for idx, part := range parts {
			hextects := strings.Split(part, ":")
			if idx == 1 {
				numColonSoFar := int(byteCount / 5)
				zeroFill := 7 - numColonSoFar - (len(hextects) - 1)
				for i := 1; i <= zeroFill*4; i++ {
					bytes[byteCount] = '0'
					byteCount++
					if i%4 == 0 {
						bytes[byteCount] = ':'
						byteCount++
					}
				}
			}
			for _, h := range hextects {
				l := len(h)
				zeroPrefix := 4 - l
				for i := 0; i < zeroPrefix; i++ {
					bytes[byteCount] = '0'
					byteCount++
				}

				for i := 0; i < l; i++ {
					bytes[byteCount] = h[i]
					byteCount++
				}
				if byteCount/5 < 7 {
					bytes[byteCount] = ':'
					byteCount++
				}
			}
		}

	}

	return IPv6Address(bytes), nil
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
