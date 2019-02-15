package asnlookup

import (
	"testing"
)

func TestParseIPv6(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want string
		err  error
	}{
		{
			name: "3 Compressed Hextect IPv6 Address Lowercase",
			ip:   "2001:db8:0:b::1a",
			want: "2001:0db8:0000:000b:0000:0000:0000:001a",
			err:  nil,
		},
		{
			name: "2 Compressed Hextect IPv6 Address Lowercase",
			ip:   "2001:db8:0:b::2a:1a",
			want: "2001:0db8:0000:000b:0000:0000:002a:001a",
			err:  nil,
		},
		/*{
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
		got, err := ParseIPv6(testCase.ip)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}
		if string(got) != testCase.want {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, string(got), testCase.want)
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
