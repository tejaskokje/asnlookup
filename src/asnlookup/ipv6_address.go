package asnlookup

import (
	"encoding/hex"
	"errors"
	//"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidIPv6Address = errors.New("Invalid IPv6 Address")
)

type IPv6Address [16]byte

func ParseIPv6(ipStr string) (IPv6Address, error) {

	ipv6Address := IPv6Address{}

	bytes := make([]byte, 32)
	ipStr = strings.ToLower(ipStr)
	if IsValidIPv6(ipStr) == false {
		return IPv6Address{}, ErrInvalidIPv6Address
	}

	byteCount := 0

	// If IPv6 address is in compressed format, convert it into
	// uncompressed canonical format without "::" or ":"
	if strings.Contains(ipStr, "::") {
		parts := strings.Split(ipStr, "::")

		for idx, part := range parts {
			hextects := strings.Split(part, ":")
			if idx == 1 {
				zeroFill := 32 - byteCount - (len(hextects) * 4)
				for i := 1; i <= zeroFill; i++ {
					bytes[byteCount] = '0'
					byteCount++
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
			}
		}
	} else {
		// TODO Handle case where we have uncompressed format but missing 0
		parts := strings.Split(ipStr, ":")

		for _, part := range parts {
			l := len(part)
			zeroFill := 4 - len(part)
			for i := 0; i < zeroFill; i++ {
				bytes[byteCount] = '0'
				byteCount++
			}

			for i := 0; i < l; i++ {
				bytes[byteCount] = part[i]
				byteCount++
			}
		}
	}

	byteCount = 0
	// We have address in uncompressed canonical form in "bytes" slice
	for i := 0; i < len(bytes); i += 2 {
		src := bytes[i : i+2]
		dst := make([]byte, 2)
		_, err := hex.Decode(dst, src)
		if err != nil {
			return ipv6Address, err
		}

		ipv6Address[byteCount] = dst[0]
		byteCount++
	}
	return ipv6Address, nil
}

func IPv6ToInt(s string) ([2]uint64, error) {

	var ipv6Int [2]uint64

	ipv6Addr, err := ParseIPv6(s)
	if err != nil {
		return ipv6Int, err
	}

	for i := 0; i < 7; i++ {
		ipv6Int[0] = (ipv6Int[0] | uint64(ipv6Addr[i])) << 8
	}
	ipv6Int[0] = (ipv6Int[0] | uint64(ipv6Addr[7]))

	for i := 8; i < 15; i++ {
		ipv6Int[1] = (ipv6Int[1] | uint64(ipv6Addr[i])) << 8
	}

	ipv6Int[1] = (ipv6Int[1] | uint64(ipv6Addr[15]))
	//fmt.Println(ipv6Int)
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
		if len(leftHextets)+len(rightHextets) > 7 {
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
