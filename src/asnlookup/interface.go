package asnlookup

// IPAddress interface contains generalized methods to abstract
// IPv4 and IPv6 addresses. Trie Insert() & Find() functions are
// IP type agnostic.
// IPv4Address & IPv6Address struct satisfy this interface.
type IPAddress interface {
	GetString() string
	GetNthHighestBit(n uint8) uint8
	GetAsn() int
	GetCidrLen() int
	GetNumBitsInAddress() int
}
