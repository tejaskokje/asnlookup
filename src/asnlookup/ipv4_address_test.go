package asnlookup

import (
	"reflect"
	"testing"
)

func TestIPv4AddressGetString(t *testing.T) {
	testCases := []struct {
		name    string
		ipCidr  string
		asn     int
		wantStr string
		err     error
	}{
		{
			name:    "Parse Valid IPv4 CIDR",
			ipCidr:  "192.168.1.0/24",
			asn:     350,
			wantStr: "192.168.1.0",
			err:     nil,
		},
		{
			name:    "Parse Incorrect CIDR Address That Can Be Corrected",
			ipCidr:  "10.10.33.4/19",
			asn:     350,
			wantStr: "10.10.32.0",
			err:     nil,
		},
		{
			name:    "Parse Incorrect CIDR Address",
			ipCidr:  "10.10.33.4/33",
			asn:     350,
			wantStr: "",
			err:     ErrInvalidIPv4Cidr,
		},
	}

	for _, testCase := range testCases {
		ipv4Address, err := NewIPv4Address(testCase.ipCidr, testCase.asn)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if err == nil {
			gotStr := ipv4Address.GetString()
			if gotStr != testCase.wantStr {
				t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, gotStr, testCase.wantStr)
			}
		}
	}
}

func TestIPv4AddressGetNthHighestBit(t *testing.T) {
	testCases := []struct {
		name   string
		ipCidr string
		asn    int
		nthBit uint8
		want   uint8
		err    error
	}{
		{
			name:   "Get 5th Bit in IPv4 Address",
			ipCidr: "192.168.1.0/24",
			asn:    350,
			nthBit: 5,
			want:   0,
			err:    nil,
		},
		{
			name:   "Get 8th Bit in IPv4 Address",
			ipCidr: "255.168.1.0/24",
			asn:    350,
			nthBit: 8,
			want:   1,
			err:    nil,
		},
		{
			name:   "Get 11th Bit in IPv4 Address",
			ipCidr: "255.168.1.0/24",
			asn:    350,
			nthBit: 11,
			want:   1,
			err:    nil,
		},
		{
			name:   "Get 24th Bit in IPv4 Address",
			ipCidr: "255.168.1.0/24",
			asn:    350,
			nthBit: 24,
			want:   1,
			err:    nil,
		},
		{
			name:   "Get 30th Bit in IPv4 Address",
			ipCidr: "255.168.1.1/30",
			asn:    350,
			nthBit: 30,
			want:   0,
			err:    nil,
		},
	}

	for _, testCase := range testCases {
		ipv4Address, err := NewIPv4Address(testCase.ipCidr, testCase.asn)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if err == nil {
			got := ipv4Address.GetNthHighestBit(testCase.nthBit)
			if got != testCase.want {
				t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
			}
		}
	}
}

func TestIPv4AddressGetAsn(t *testing.T) {
	testCases := []struct {
		name   string
		ipCidr string
		asn    int
		want   int
		err    error
	}{
		{
			name:   "Parse Valid IPv4 CIDR",
			ipCidr: "192.168.1.0/24",
			asn:    350,
			want:   350,
			err:    nil,
		},
		{
			name:   "Parse Incorrect CIDR Address That Can Be Corrected",
			ipCidr: "10.10.33.4/19",
			asn:    351,
			want:   351,
			err:    nil,
		},
		{
			name:   "Parse Incorrect CIDR Address",
			ipCidr: "10.10.33.4/33",
			asn:    352,
			want:   352,
			err:    ErrInvalidIPv4Cidr,
		},
	}

	for _, testCase := range testCases {
		ipv4Address, err := NewIPv4Address(testCase.ipCidr, testCase.asn)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if err == nil {
			got := ipv4Address.GetAsn()
			if got != testCase.want {
				t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
			}
		}
	}
}

func TestIPv4AddressGetCidrLen(t *testing.T) {
	testCases := []struct {
		name   string
		ipCidr string
		asn    int
		want   int
		err    error
	}{
		{
			name:   "Parse Valid IPv4 CIDR",
			ipCidr: "192.168.1.0/24",
			asn:    350,
			want:   24,
			err:    nil,
		},
		{
			name:   "Parse Incorrect CIDR Address That Can Be Corrected",
			ipCidr: "10.10.33.4/19",
			asn:    351,
			want:   19,
			err:    nil,
		},
	}

	for _, testCase := range testCases {
		ipv4Address, err := NewIPv4Address(testCase.ipCidr, testCase.asn)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if err == nil {
			got := ipv4Address.GetCidrLen()
			if got != testCase.want {
				t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
			}
		}
	}
}

func TestIPv4AddressGetNumBitsInAddress(t *testing.T) {
	testCases := []struct {
		name   string
		ipCidr string
		asn    int
		want   int
		err    error
	}{
		{
			name:   "Parse Valid IPv4 CIDR",
			ipCidr: "192.168.1.0/24",
			asn:    350,
			want:   32,
			err:    nil,
		},
		{
			name:   "Parse Incorrect CIDR Address That Can Be Corrected",
			ipCidr: "10.10.33.4/19",
			asn:    351,
			want:   32,
			err:    nil,
		},
	}

	for _, testCase := range testCases {
		ipv4Address, err := NewIPv4Address(testCase.ipCidr, testCase.asn)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if err == nil {
			got := ipv4Address.GetNumBitsInAddress()
			if got != testCase.want {
				t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
			}
		}
	}
}

func TestParseIPv4(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want []byte
		err  error
	}{
		{
			name: "Parse Valid IPv4 Address",
			ip:   "192.168.1.1",
			want: []byte{192, 168, 1, 1},
			err:  nil,
		},
		{
			name: "Parse Invalid IPv4 Address",
			ip:   "500.500.500.1",
			want: []byte{},
			err:  ErrInvalidIPv4Address,
		},
		{
			name: "Parse Invalid IPv4 Address With All 0",
			ip:   "0.0.0.0",
			want: []byte{},
			err:  ErrInvalidIPv4Address,
		},
		{
			name: "Parse Valid IPv4 Address With 0",
			ip:   "192.0.0.1",
			want: []byte{192, 0, 0, 1},
			err:  nil,
		},
	}

	for _, testCase := range testCases {
		got, err := parseIPv4(testCase.ip)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if reflect.DeepEqual(got, testCase.want) != true {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}

func TestIPv4StrToInt(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want uint32
		err  error
	}{
		{
			name: "Valid IPv4 Address",
			ip:   "8.8.8.8",
			want: 134744072,
			err:  nil,
		},
		{
			name: "Valid IPv4 Address With All 1s",
			ip:   "255.255.255.255",
			want: 4294967295,
			err:  nil,
		},
		{
			name: "InValid IPv4 Address With All 0s",
			ip:   "0.0.0.0",
			want: 0,
			err:  ErrInvalidIPv4Address,
		},
		{
			name: "Valid IPv4 Address With 0s In Lower Three Octets",
			ip:   "254.0.0.0",
			want: 4261412864,
			err:  nil,
		},
	}

	for _, testCase := range testCases {
		got, err := ipv4StrToInt(testCase.ip)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if got != testCase.want {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}

func TestIntToIPv4Str(t *testing.T) {
	testCases := []struct {
		name string
		ip   uint32
		want string
		err  error
	}{
		{
			name: "Valid IPv4 Address",
			ip:   134744072,
			want: "8.8.8.8",
			err:  nil,
		},
		{
			name: "Valid IPv4 Address With 0 In Middle Octets",
			ip:   3221225473,
			want: "192.0.0.1",
			err:  nil,
		},
		{
			name: "Valid IPv4 Address With All 1s",
			ip:   4294967295,
			want: "255.255.255.255",
			err:  nil,
		},
		{
			name: "InValid IPv4 Address With All 0s",
			ip:   0,
			want: "",
			err:  ErrInvalidIPv4Address,
		},
		{
			name: "Valid IPv4 Address With 0s In Lower Three Octets",
			ip:   4261412864,
			want: "254.0.0.0",
			err:  nil,
		},
	}

	for _, testCase := range testCases {
		got, err := intToIPv4Str(testCase.ip)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if got != testCase.want {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}

func TestIsValidIPv4(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "Valid IPv4 Address",
			ip:   "8.8.8.8",
			want: true,
		},
		{
			name: "Valid IPv4 Address With 0 in Second, Third and Fourth Octet",
			ip:   "8.0.0.0",
			want: true,
		},
		{
			name: "Invalid IPv4 Address With All 0",
			ip:   "0.0.0.0",
			want: false,
		},
		{
			name: "Invalid IPv4 Address With 0 In Highest Octet",
			ip:   "0.1.1.1",
			want: false,
		},
		{
			name: "Not Enough Octets for IPv4 Address",
			ip:   "192.168.1.",
			want: false,
		},
		{
			name: "Invalid High Octet range for IPv4 Address",
			ip:   "300.168.1.1",
			want: false,
		},
		{
			name: "Invalid All Octets for IPv4 Address ",
			ip:   "300.300.500.500",
			want: false,
		},
		{
			name: "More Than 3 Dots In IPv4 Address",
			ip:   "192.168.1.1.1",
			want: false,
		},
		{
			name: "IPv4 Address With 0 As Last Octet",
			ip:   "8.8.8.0",
			want: true,
		},

		{
			name: "Invalid Characters In IPv4 Address",
			ip:   "192-168.1.1.1",
			want: false,
		},
	}

	for _, testCase := range testCases {
		got := IsValidIPv4(testCase.ip)
		if got != testCase.want {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}

func TestIsValidIPv4Cidr(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "Valid CIDR Host",
			ip:   "192.168.1.1/32",
			want: true,
		},
		{
			name: "Valid CIDR Subnet",
			ip:   "192.168.1.0/24",
			want: true,
		},
		{
			name: "Invalid CIDR For IPv4 Address",
			ip:   "8.8.8.8/33",
			want: false,
		},
		{
			name: "Zero CIDR for IPv4 Address",
			ip:   "5.5.5.5/0",
			want: false,
		},
		{
			name: "More Than One / In Mask At Different Places",
			ip:   "192.168./1.1/25",
			want: false,
		},
	}

	for _, testCase := range testCases {
		got := IsValidIPv4Cidr(testCase.ip)
		if got != testCase.want {
			t.Errorf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}
