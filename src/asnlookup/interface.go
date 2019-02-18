package asnlookup

type IPAddress interface {
	GetString() string
	GetNthHighestBit(n uint8) uint8
	GetAsn() int
	GetCidrLen() int
	GetNumBitsInAddress() int
	//DumpAddressInBinary() string
}
