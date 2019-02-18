package asnlookup

import (
	"reflect"
	"testing"
)

type ipCidrAsn struct {
	ip  string
	asn int
}

func TestInsertFindIPv4(t *testing.T) {
	testCases := []struct {
		name       string
		ipCidrList []ipCidrAsn
		ipToFind   string
		want       NodeInfoList
		err        error
	}{
		{
			name: "Insert-Find IPv4 Address",
			ipCidrList: []ipCidrAsn{
				{
					"192.168.1.1/24",
					351,
				},
				{
					"192.168.0.0/16",
					355,
				},
				{
					"8.8.8.0/14",
					450,
				},
				{
					"192.168.1.0/20",
					600,
				},
				{
					"1.1.1.0/16",
					200,
				},
			},

			ipToFind: "192.168.1.5/32",
			want: NodeInfoList{
				{"192.168.1.0", 24, 351},
				{"192.168.0.0", 20, 600},
				{"192.168.0.0", 16, 355},
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		trie := NewTrie()
		for _, ipCidr := range testCase.ipCidrList {
			ipv4Address, err := NewIPv4Address(ipCidr.ip, ipCidr.asn)
			if err != testCase.err {
				t.Fatalf("%s: received error for %s/%d does not match: got %v, want %v", testCase.name, ipCidr.ip, ipCidr.asn, err, testCase.err)
			}

			Insert(trie, ipv4Address)
		}

		ipv4Address, err := NewIPv4Address(testCase.ipToFind, -1)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		got := Find(trie, ipv4Address)
		if reflect.DeepEqual(got, testCase.want) != true {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}

}

func TestInsertFindIPv6(t *testing.T) {
	testCases := []struct {
		name       string
		ipCidrList []ipCidrAsn
		ipToFind   string
		want       NodeInfoList
		err        error
	}{
		{
			name: "Insert-Find IPv6 Address",
			ipCidrList: []ipCidrAsn{
				{
					"2001:db8:0:b::1A:1c/64",
					451,
				},
				{
					"2604:a880:2:d0::2249:2001/77",
					455,
				},
				{
					"2001:db8:0:b::1A:1c/67",
					550,
				},
				{
					"fe80::c7e:afff:fe10:66e0/64",
					700,
				},
				{
					"2001:db8:0:b::1A:1c/80",
					300,
				},
			},

			ipToFind: "2001:db8:0:b::1A:1c/128",
			want: NodeInfoList{
				{"2001:0db8:0000:000b:0000:0000:0000:0000", 80, 300},
				{"2001:0db8:0000:000b:0000:0000:0000:0000", 67, 550},
				{"2001:0db8:0000:000b:0000:0000:0000:0000", 64, 451},
			},
			err: nil,
		},
	}

	for _, testCase := range testCases {
		trie := NewTrie()
		for _, ipCidr := range testCase.ipCidrList {
			ipv6Address, err := NewIPv6Address(ipCidr.ip, ipCidr.asn)
			if err != testCase.err {
				t.Fatalf("%s: received error for %s/%d does not match: got %v, want %v", testCase.name, ipCidr.ip, ipCidr.asn, err, testCase.err)
			}

			Insert(trie, ipv6Address)
		}

		ipv6Address, err := NewIPv6Address(testCase.ipToFind, -1)
		if err != testCase.err {
			t.Fatalf("%s: received error does not match: got %v, want %v", testCase.name, err, testCase.err)
		}

		got := Find(trie, ipv6Address)
		if reflect.DeepEqual(got, testCase.want) != true {
			t.Fatalf("%s: result does not match: got %v, want %v", testCase.name, got, testCase.want)
		}
	}
}
