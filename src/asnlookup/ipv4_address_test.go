package asnlookup

import (
	"testing"
)

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
