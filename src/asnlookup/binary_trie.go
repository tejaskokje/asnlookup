package asnlookup

import (
	"fmt"
	"sort"
)

type NodeInfo struct {
	Subnet string
	Cidr   int
	Asn    int
}

type NodeInfoList []NodeInfo

func (n NodeInfoList) Len() int {
	return len(n)
}

func (n NodeInfoList) Less(i, j int) bool {
	if n[i].Cidr > n[j].Cidr {
		return true
	}
	return false
}

func (n NodeInfoList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

type Node struct {
	Info  []NodeInfo
	Left  *Node
	Right *Node
}

type Trie struct {
	TrieType int
	Size     int
	Root     *Node
}

func NewTrie() *Trie {
	return &Trie{
		TrieType: 0, //v4
		Root:     NewNode(),
	}
}

func NewNode() *Node {
	return &Node{
		Info:  []NodeInfo{},
		Left:  nil,
		Right: nil,
	}

}

func Insert(t *Trie, ip IPAddress) {
	// Safe to ignore error below as key will already be sanitized by this time
	root := t.Root
	for i := 1; i <= ip.GetCidrLen(); i++ {
		child := ip.GetNthHighestBit(uint8(i))
		if child == 0 {
			if root.Left == nil {
				root.Left = NewNode()
			}

			root = root.Left
		} else {
			if root.Right == nil {
				root.Right = NewNode()
			}

			root = root.Right
		}
	}

	root.Info = append(root.Info, NodeInfo{ip.GetString(), ip.GetCidrLen(), ip.GetAsn()})
	return
}

func Find(cfg *Config) NodeInfoList {
	infoList := NodeInfoList{}
	root := cfg.trie.Root
	for i := 1; i <= cfg.IPToFind.GetNumBitsInAddress(); i++ {
		child := cfg.IPToFind.GetNthHighestBit(uint8(i))
		if child == 0 && root.Left != nil {
			if len(root.Left.Info) > 0 {
				infoList = append(infoList, root.Left.Info...)
			}
			root = root.Left
		} else if child == 1 && root.Right != nil {

			if len(root.Right.Info) > 0 {
				infoList = append(infoList, root.Right.Info...)
			}

			root = root.Right
		} else {
			break
		}
	}

	sort.Sort(infoList)
	return infoList
}

func DumpTrie(t *Trie) {
	root := t.Root
	if root != nil {
		DumpNode(root, -1)
	}
}

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
