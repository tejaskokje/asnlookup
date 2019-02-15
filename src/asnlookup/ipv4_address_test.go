package asnlookup

import (
	"testing"
)

func TestParseIPv4(t *testing.T) {
	testCases := []struct {
		name string
		ip   string
		want IPv4Address
		err  error
	}{
		{
			name: "Parse Valid IPv4 Address",
			ip:   "192.168.1.1",
			want: IPv4Address{192, 168, 1, 1},
			err:  nil,
		},
		{
			name: "Parse Invalid IPv4 Address",
			ip:   "500.500.500.1",
			want: IPv4Address{},
			err:  ErrInvalidIPv4Address,
		},
		{
			name: "Parse Invalid IPv4 Address With All 0",
			ip:   "0.0.0.0",
			want: IPv4Address{},
			err:  ErrInvalidIPv4Address,
		},
		{
			name: "Parse Valid IPv4 Address With 0",
			ip:   "192.0.0.1",
			want: IPv4Address{192, 0, 0, 1},
			err:  nil,
		},
	}

	for _, testCase := range testCases {
		got, err := ParseIPv4(testCase.ip)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		if got != testCase.want {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}

func TestIPv4ToInt(t *testing.T) {
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
		got, err := IPv4ToInt(testCase.ip)
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
