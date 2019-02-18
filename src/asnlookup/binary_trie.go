package asnlookup

import (
	"fmt"
	"sort"
)

// NodeInfo stores subnet, Cidr and Asn information.
// This structure is agnostic to IPv4 or IPv6 address
type NodeInfo struct {
	Subnet string
	Cidr   int
	Asn    int
}

// NodeInfoList stores a list of NodeInfo in a given trie node
// This type is also used for sort interface to sort according
// to Cidr for displaying to user
type NodeInfoList []NodeInfo

// Len implements Len() method for sort interface
func (n NodeInfoList) Len() int {
	return len(n)
}

// Less implements Less() method for sort interface
func (n NodeInfoList) Less(i, j int) bool {
	if n[i].Cidr > n[j].Cidr {
		return true
	}
	return false
}

// Swap implements Swap() method for sort interface
func (n NodeInfoList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// Node is a trie node
type Node struct {
	Info  []NodeInfo
	Left  *Node
	Right *Node
}

// Trie struct holds information about trie.
// Currently it only holds root of trie. It can be
// enhanced to store size & content info (IPv4 or IPv6)
type Trie struct {
	Root *Node
}

// NewTrie creates a Trie and returns its pointer
func NewTrie() *Trie {
	return &Trie{
		Root: NewNode(),
	}
}

// NewNode creates a new trie node
func NewNode() *Node {
	return &Node{
		Info:  []NodeInfo{},
		Left:  nil,
		Right: nil,
	}

}

// Insert adds a node into the trie. Input "ip" can either be IPv4 or IPv6 address.
// Trie is agnostic to IP address type as it works on 0s and 1s.
// IPv4 trie can have maximum 32 lookups. IPv6 trie can have 128 lookups.
func Insert(t *Trie, ip IPAddress) {
	// Safe to ignore error below as key will already be sanitized by this time
	root := t.Root

	// Get the Cidr prefix length and iterate over bits starting with highest
	// order bit.
	for i := 1; i <= ip.GetCidrLen(); i++ {
		child := ip.GetNthHighestBit(uint8(i))

		// 0 bit means go to left child, 1 bit means go to right child
		if child == 0 {
			// If left node is nil, create a new left trie node
			if root.Left == nil {
				root.Left = NewNode()
			}

			root = root.Left
		} else {
			// If right node is nil, create a new right trie node
			if root.Right == nil {
				root.Right = NewNode()
			}

			root = root.Right
		}
	}

	// We are done interating over all bits of Cidr prefix. Store information in NodeList
	// for current trie node
	root.Info = append(root.Info, NodeInfo{ip.GetString(), ip.GetCidrLen(), ip.GetAsn()})
	return
}

// Find walks through the bits of target IP address and returns NodeInfoList
// with matching trie nodes for target IP address
func Find(cfg *Config) NodeInfoList {
	infoList := NodeInfoList{}
	root := cfg.trie.Root

	for i := 1; i <= cfg.IPToFind.GetNumBitsInAddress(); i++ {
		child := cfg.IPToFind.GetNthHighestBit(uint8(i))
		if child == 0 && root.Left != nil {
			// Left child matches with target IP. Store it in infoList
			if len(root.Left.Info) > 0 {
				infoList = append(infoList, root.Left.Info...)
			}
			root = root.Left
		} else if child == 1 && root.Right != nil {
			// Right child matches with target IP. Store it in infoList
			if len(root.Right.Info) > 0 {
				infoList = append(infoList, root.Right.Info...)
			}

			root = root.Right
		} else {
			break
		}
	}

	// Return sorted infoList by Cidr length
	sort.Sort(infoList)
	return infoList
}

// DumpTrie dumps trie for debugging
func DumpTrie(t *Trie) {
	root := t.Root
	if root != nil {
		DumpNode(root, -1)
	}
}

// DumpNode dumps trie nodes for debugging
func DumpNode(n *Node, dir int) {
	relation := ""
	if n != nil {
		switch dir {
		case -1:
			relation = "root"
		case 0:
			relation = "left"
		case 1:
			relation = "right"
		default:
			relation = "unknown"
		}

		fmt.Println(relation)

		for _, info := range n.Info {
			fmt.Println("\t", info.Subnet, info.Cidr, info.Asn)
		}

		DumpNode(n.Left, 0)
		DumpNode(n.Right, 1)
	}
}
