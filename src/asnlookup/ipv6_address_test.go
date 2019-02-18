package asnlookup

import (
	"reflect"
	"testing"
)

func TestIPv6AddressGetString(t *testing.T) {
	testCases := []struct {
		name    string
		ipCidr  string
		asn     int
		wantStr string
		err     error
	}{
		{
			name:    "Parse Valid IPv6 CIDR",
			ipCidr:  "2604:a880:2:d0::2249:2001/64",
			asn:     350,
			wantStr: "2604:a880:0002:00d0:0000:0000:0000:0000",
			err:     nil,
		},
		{
			name:    "Parse Incorrect CIDR Address That Can Be Corrected",
			ipCidr:  "fe80::c7e:afff:fe10:66e0:12/64",
			asn:     350,
			wantStr: "fe80:0000:0000:0c7e:0000:0000:0000:0000",
			err:     nil,
		},
		{
			name:    "Parse Incorrect CIDR Address",
			ipCidr:  "fe80::c7e:afff:fe10:66e0/129",
			asn:     350,
			wantStr: "",
			err:     ErrInvalidIPv6Cidr,
		},
	}

	for _, testCase := range testCases {
		ipv4Address, err := newIPv6Address(testCase.ipCidr, testCase.asn)
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

func TestIPv6AddressGetNthHighestBit(t *testing.T) {
	testCases := []struct {
		name   string
		ipCidr string
		asn    int
		nthBit uint8
		want   uint8
		err    error
	}{
		{
			name:   "Get 9th Bit in IPv6 Address",
			ipCidr: "fe80::c7e:afff:fe10:66e0:12/64",
			asn:    350,
			nthBit: 9,
			want:   1,
			err:    nil,
		},
		{
			name:   "Get 25nd Bit in IPv6 Address",
			ipCidr: "fe80::c7e:afff:fe10:66e0:12/64",
			asn:    350,
			nthBit: 25,
			want:   0,
			err:    nil,
		},
		{
			name:   "Get 33rd Bit in IPv4 Address",
			ipCidr: "fe80::c7e:afff:fe10:66e0:12/64",
			asn:    350,
			nthBit: 33,
			want:   0,
			err:    nil,
		},
		{
			name:   "Get 51st Bit in IPv4 Address",
			ipCidr: "fe80::ffff:ffff:ffff:66e0:12/128",
			asn:    350,
			nthBit: 51,
			want:   1,
			err:    nil,
		},
		{
			name:   "Get 77th Bit in IPv4 Address",
			ipCidr: "fe80::ffff:ffff:ffff:66e0:12/128",
			asn:    350,
			nthBit: 77,
			want:   1,
			err:    nil,
		},
		{
			name:   "Get 113th Bit in IPv4 Address",
			ipCidr: "fe80::ffff:ffff:ffff:66e0:12/128",
			asn:    350,
			nthBit: 113,
			want:   0,
			err:    nil,
		},
	}

	for _, testCase := range testCases {
		ipv4Address, err := newIPv6Address(testCase.ipCidr, testCase.asn)
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
func TestIPv6AddressGetAsn(t *testing.T) {
	testCases := []struct {
		name   string
		ipCidr string
		asn    int
		want   int
		err    error
	}{
		{
			name:   "Parse Valid IPv6 CIDR",
			ipCidr: "2604:a880:2:d0::2249:2001/64",
			asn:    350,
			want:   350,
			err:    nil,
		},
		{
			name:   "Parse Incorrect CIDR Address That Can Be Corrected",
			ipCidr: "fe80::c7e:afff:fe10:66e0:12/64",
			asn:    351,
			want:   351,
			err:    nil,
		},
	}

	for _, testCase := range testCases {
		ipv4Address, err := newIPv6Address(testCase.ipCidr, testCase.asn)
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

func TestIPv6AddressGetCidrLen(t *testing.T) {
	testCases := []struct {
		name   string
		ipCidr string
		asn    int
		want   int
		err    error
	}{
		{
			name:   "Parse Valid IPv6 CIDR",
			ipCidr: "2604:a880:2:d0::2249:2001/64",
			asn:    350,
			want:   64,
			err:    nil,
		},
		{
			name:   "Parse Incorrect CIDR Address That Can Be Corrected",
			ipCidr: "fe80::c7e:afff:fe10:66e0:12/77",
			asn:    351,
			want:   77,
			err:    nil,
		},
		{
			name:   "Parse Incorrect CIDR Address That Can Be Corrected",
			ipCidr: "fe80::c7e:afff:fe10:66e0:12/128",
			asn:    351,
			want:   128,
			err:    nil,
		},
	}

	for _, testCase := range testCases {
		ipv4Address, err := newIPv6Address(testCase.ipCidr, testCase.asn)
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

func TestIPv6AddressGetNumBitsInAddress(t *testing.T) {
	testCases := []struct {
		name   string
		ipCidr string
		asn    int
		want   int
		err    error
	}{
		{
			name:   "Parse Valid IPv6 CIDR",
			ipCidr: "2604:a880:2:d0::2249:2001/64",
			asn:    350,
			want:   128,
			err:    nil,
		},
		{
			name:   "Parse Incorrect CIDR Address That Can Be Corrected",
			ipCidr: "fe80::c7e:afff:fe10:66e0:12/77",
			asn:    351,
			want:   128,
			err:    nil,
		},
	}

	for _, testCase := range testCases {
		ipv4Address, err := newIPv6Address(testCase.ipCidr, testCase.asn)
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
func TestParseIPv6(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want []byte
		err  error
	}{
		{
			name: "3 Compressed Hextect IPv6 Address Lowercase",
			ip:   "2001:db8:0:b::1a",
			want: []byte{32, 1, 13, 184, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0, 26},
			err:  nil,
		},
		{
			name: "2 Compressed Hextect IPv6 Address Lowercase",
			ip:   "2001:db8:0:b::2a:1a",
			want: []byte{32, 1, 13, 184, 0, 0, 0, 11, 0, 0, 0, 0, 0, 42, 0, 26},
			err:  nil,
		},
		{
			name: "Uncompressed Hextect IPv6 Address Lowercase",
			ip:   "2001:0db8:0000:000b:0000:0000:0000:001a",
			want: []byte{32, 1, 13, 184, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0, 26},
			err:  nil,
		},
		{
			name: "Error In IPv6 Address Format",
			ip:   "2001:0db8:0000:000b:0000:0000:0000::001a",
			want: []byte{},
			err:  ErrInvalidIPv6Address,
		},
		{
			name: "Uncompressed Hextect IPv6 Address Lowercase With Missing 0s",
			ip:   "2001:db8:0:b:00:00:0:1a",
			want: []byte{32, 1, 13, 184, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0, 26},
			err:  nil,
		},
	}

	for _, testCase := range testCases {
		got, err := parseIPv6(testCase.ip)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if reflect.DeepEqual(got, testCase.want) != true {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}

func TestIsValidIPv6(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "Compressed IPv6 Address Lowercase",
			ip:   "2001:db8:0:b::1a",
			want: true,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase",
			ip:   "2001:0db8:0000:000b:0000:0000:0000:001a",
			want: true,
		},
		{
			name: "More than one compression IPv6 Address Lowercase",
			ip:   "2001:db8::0:b::1a",
			want: false,
		},

		{
			name: "Compressed IPv6 Address Uppercase",
			ip:   "2001:DB8:0:b::1A",
			want: true,
		},
		{
			name: "Uncompressed IPv6 Address Uppercase",
			ip:   "2001:0DB8:0000:000B:0000:0000:0000:001A",
			want: true,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase multiple ::",
			ip:   "2001:0DB8:0000:000B:0000::0000:0000::001A",
			want: false,
		},
		{
			name: "More than one compression IPv6 Address Uppercase",
			ip:   "2001:DB8::0:b::1A",
			want: false,
		},

		{
			name: "Compressed IPv6 Address Mixedcase",
			ip:   "2001:db8:C:b::1A",
			want: true,
		},
		{
			name: "Uncompressed IPv6 Address Mixedcase",
			ip:   "2001:0db8:0000:000b:0C00:0000:0000:001A",
			want: true,
		},
		{
			name: "More than one compression IPv6 Address Mixedcase",
			ip:   "2001:db8::C:b::1A",
			want: false,
		},
	}

	for _, testCase := range testCases {
		got := isValidIPv6(testCase.ip)
		if got != testCase.want {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}

func TestIPv6StrToInt(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want [2]uint64
		err  error
	}{
		{
			name: "Compressed IPv6 Address Lowercase",
			ip:   "2001:db8:0:b::1a",
			want: [2]uint64{2306139568115548171, 26},
			err:  nil,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase",
			ip:   "2001:0db8:0000:000b:0000:02b0:0000:001a",
			want: [2]uint64{2306139568115548171, 2954937499674},
			err:  nil,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase With 0 In Higher 64 Bits",
			ip:   "0000:0000:0000:0000:0000:02b0:0000:001a",
			want: [2]uint64{0, 2954937499674},
			err:  nil,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase With 0 In Lower 64 Bits",
			ip:   "2001:0db8:0000:000b:0000:0000:0000:0000",
			want: [2]uint64{2306139568115548171, 0},
			err:  nil,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase With 1 In Higher 64 Bits",
			ip:   "ffff:ffff:ffff:ffff:0000:02b0:0000:001a",
			want: [2]uint64{18446744073709551615, 2954937499674},
			err:  nil,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase With 1 In Lower 64 Bits",
			ip:   "2001:0db8:0000:000b:ffff:ffff:ffff:ffff",
			want: [2]uint64{2306139568115548171, 18446744073709551615},
			err:  nil,
		},
	}

	for _, testCase := range testCases {
		got, err := ipv6StrToInt(testCase.ip)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if got != testCase.want {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}
func TestIntToIPv6Str(t *testing.T) {
	testCases := []struct {
		name string
		ip   [2]uint64
		want string
		err  error
	}{
		{
			name: "Compressed IPv6 Address Lowercase",
			ip:   [2]uint64{2306139568115548171, 26},
			want: "2001:0db8:0000:000b:0000:0000:0000:001a",
			err:  nil,
		},

		{
			name: "Uncompressed IPv6 Address Lowercase",
			ip:   [2]uint64{2306139568115548171, 2954937499674},
			want: "2001:0db8:0000:000b:0000:02b0:0000:001a",
			err:  nil,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase With 0 In Higher 64 Bits",
			ip:   [2]uint64{0, 2954937499674},
			want: "0000:0000:0000:0000:0000:02b0:0000:001a",
			err:  nil,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase With 0 In Lower 64 Bits",
			ip:   [2]uint64{2306139568115548171, 0},
			want: "2001:0db8:0000:000b:0000:0000:0000:0000",
			err:  nil,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase With 1 In Higher 64 Bits",
			ip:   [2]uint64{18446744073709551615, 2954937499674},
			want: "ffff:ffff:ffff:ffff:0000:02b0:0000:001a",
			err:  nil,
		},
		{
			name: "Uncompressed IPv6 Address Lowercase With 1 In Lower 64 Bits",
			ip:   [2]uint64{2306139568115548171, 18446744073709551615},
			want: "2001:0db8:0000:000b:ffff:ffff:ffff:ffff",
			err:  nil,
		},
	}

	for _, testCase := range testCases {
		got, err := intToIPv6Str(testCase.ip)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if got != testCase.want {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}
func TestIsValidIPv6Cidr(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "Valid CIDR",
			ip:   "2001:db8:0:b::1a/128",
			want: true,
		},
		{
			name: "Uncompressed IPv6 Address Valid CIDR",
			ip:   "2001:0db8:0000:000b:0000:0000:0000:001a/128",
			want: true,
		},
		{
			name: "IPv6 Mask More Than 128",
			ip:   "2001:db8::0:b:1a/129",
			want: false,
		},
		{
			name: "More than one / in mask",
			ip:   "2001:db8::0:b::1a//129",
			want: false,
		},
		{
			name: "More than one / in mask at different places",
			ip:   "2001:db8::0:/b::1a/129",
			want: false,
		},
	}

	for _, testCase := range testCases {
		got := isValidIPv6Cidr(testCase.ip)
		if got != testCase.want {
			t.Errorf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}
