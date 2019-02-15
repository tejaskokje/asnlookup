package asnlookup

import (
	"testing"
)

func TestParseIPv6(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want IPv6Address
		err  error
	}{
		{
			name: "3 Compressed Hextect IPv6 Address Lowercase",
			ip:   "2001:db8:0:b::1a",
			want: IPv6Address{32, 1, 13, 184, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0, 26},
			err:  nil,
		},
		{
			name: "2 Compressed Hextect IPv6 Address Lowercase",
			ip:   "2001:db8:0:b::2a:1a",
			want: IPv6Address{32, 1, 13, 184, 0, 0, 0, 11, 0, 0, 0, 0, 0, 42, 0, 26},
			err:  nil,
		},
		{
			name: "Uncompressed Hextect IPv6 Address Lowercase",
			ip:   "2001:0db8:0000:000b:0000:0000:0000:001a",
			want: IPv6Address{32, 1, 13, 184, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0, 26},
			err:  nil,
		},
		{
			name: "Error In IPv6 Address Format",
			ip:   "2001:0db8:0000:000b:0000:0000:0000::001a",
			want: IPv6Address{},
			err:  ErrInvalidIPv6Address,
		},
		{
			name: "Uncompressed Hextect IPv6 Address Lowercase with missing 0s",
			ip:   "2001:db8:0000:b:0000:0000:0000:1a",
			want: IPv6Address{32, 1, 13, 184, 0, 0, 0, 11, 0, 0, 0, 0, 0, 0, 0, 26},
			err:  nil,
		},
	}

	for _, testCase := range testCases {
		got, err := ParseIPv6(testCase.ip)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if got != testCase.want {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}

func TestIPv6ToInt(t *testing.T) {
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
		/*
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
			},*/
	}

	for _, testCase := range testCases {
		got, err := IPv6ToInt(testCase.ip)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if got != testCase.want {
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
		got := IsValidIPv6(testCase.ip)
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
		got := IsValidIPv6Cidr(testCase.ip)
		if got != testCase.want {
			t.Errorf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}
