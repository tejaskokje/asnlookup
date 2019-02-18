package asnlookup

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrInvalidIPv6Address is returned when IPv6 address is badly formatted
	ErrInvalidIPv6Address = errors.New("Invalid IPv6 Address")

	// ErrInvalidIPv6Cidr is returned when IPv6 address & CIDR is badly formatted
	ErrInvalidIPv6Cidr = errors.New("Invalid IPv6 CIDR Format")
)

// IPv6Address struct satisfies IPAddress interface. This is used in trie Insert()
// & Find() functions to make them IP address type agnostic in those functions
type IPv6Address struct {
	cidrLen int
	mask    [2]uint64
	ip      [2]uint64
	ipStr   string
	asn     int
}

// Compile time check to ensure IPv6Address satiesfies IPAddress interface
var _ IPAddress = &IPv6Address{}

// newIPv4Address returns new IPv4Address
func newIPv6Address(ipCidr string, asn int) (IPAddress, error) {
	if isValidIPv6Cidr(ipCidr) == false {
		return nil, ErrInvalidIPv6Cidr
	}

	ipv6Address := IPv6Address{}
	ip := strings.Split(ipCidr, "/")
	ipInt, err := ipv6StrToInt(ip[0])
	if err != nil {
		return ipv6Address, err
	}

	prefix, err := strconv.Atoi(ip[1])
	if err != nil {
		return ipv6Address, err
	}

	ipv6Address.cidrLen = prefix

	// Adjust masks to correctly "and" them with ip address array
	mask1 := 0
	mask2 := 0
	if prefix > 64 {
		mask1 = 0
		mask2 = 128 - prefix
	} else {
		mask1 = 64 - prefix
		mask2 = 64
	}

	// Generate mask for IPv4 address
	ipv6Address.mask[0] = uint64(^(uint64(0))) << uint64(mask1)
	ipv6Address.mask[1] = uint64(^(uint64(0))) << uint64(mask2)

	// Store correct subnet even if input has host bits set in the address
	ipv6Address.ip[0] = (ipInt[0] & ipv6Address.mask[0])
	ipv6Address.ip[1] = (ipInt[1] & ipv6Address.mask[1])

	ipStr, err := intToIPv6Str(ipv6Address.ip)
	if err != nil {
		return ipv6Address, err
	}
	ipv6Address.ipStr = ipStr
	ipv6Address.asn = asn

	return ipv6Address, nil
}

// GetString returns IPv6 address in string format. This method
// is needed to satisfy IPAddress interface
func (ipv6 IPv6Address) GetString() string {
	return ipv6.ipStr
}

// GetNthHighestBit returns nth highest bit for IPv6 address. This method
// is needed to satisfy IPAddress interface
func (ipv6 IPv6Address) GetNthHighestBit(n uint8) uint8 {

	var nthBit uint64
	if n <= 64 {
		nthBit = ((ipv6.ip[0]) >> (64 - uint32(n))) & 0x1
	} else if n > 64 && n <= 128 {
		nthBit = ((ipv6.ip[1]) >> (64 - uint32(n-64))) & 0x1
	} else {
		nthBit = 255
	}

	return uint8(nthBit)
}

// GetAsn returns ASN stored in IPv6 address. This method
// is needed to satisfy IPAddress interface
func (ipv6 IPv6Address) GetAsn() int {
	return ipv6.asn
}

// GetCidrLen returns CIDR prefix length stored in IPv6 address.
// This method is needed to satisfy IPAddress interface
func (ipv6 IPv6Address) GetCidrLen() int {
	return ipv6.cidrLen
}

// GetNumBitsInAddress returns number of bits in IPv6 address.
// This method is needed to satisfy IPAddress interface
func (ipv6 IPv6Address) GetNumBitsInAddress() int {
	return 128
}

// Following are helper functions to parse, validate & convert IPv6 address
func parseIPv6(ipStr string) ([]byte, error) {

	if isValidIPv6(ipStr) == false {
		return []byte{}, ErrInvalidIPv6Address
	}

	bytes := make([]byte, 32)
	ipStr = strings.ToLower(ipStr)
	ipv6Address := make([]byte, 16)
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

func isValidIPv6(ip string) bool {

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

func ipv6StrToInt(s string) ([2]uint64, error) {

	var ipv6Int [2]uint64

	ipv6Addr, err := parseIPv6(s)
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
	return ipv6Int, nil
}

// This function will return IPv6 address string in uncompressed format.
// It can be enhanced to return compressed format string in future
func intToIPv6Str(ip [2]uint64) (string, error) {

	var ipStr []string
	for _, part := range ip {
		for i := 1; i <= 4; i++ {
			octet := (part >> uint64(64-16*i)) & 0xFFFF
			s := fmt.Sprintf("%04x", octet)
			ipStr = append(ipStr, s)
		}
	}
	return strings.Join(ipStr, ":"), nil
}

func isValidIPv6Cidr(cidr string) bool {
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return false
	}

	if isValidIPv6(parts[0]) == false {
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
