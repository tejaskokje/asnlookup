package asnlookup

import (
	"fmt"
	"sort"
)

type NodeInfo struct {
	Subnet uint32
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
		Size:     0,
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
func Insert(t *Trie, key uint32, value int, prefixLen int) {
	//fmt.Println("Inserting ", key, value, prefixLen)
	root := t.Root
	origKey := key
	//fmt.Printf("%d %032b, %b, %032b\n", key, key, 0x80000000, (key << uint32(1)))
	for i := 1; i <= 32; i++ {
		child := ((key) >> 31) & 0x1
		key = key << 1
		//fmt.Printf("%b | ", child)
		//fmt.Printf("child %b key %032b\n", child, key)

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
		if key == 0 {
			break
		}
	}

	root.Info = append(root.Info, NodeInfo{origKey, prefixLen, value})
	return
}

func Find(t *Trie, key uint32) NodeInfoList {
	//fmt.Printf("Find Key %032b\n", key)
	infoList := NodeInfoList{}
	root := t.Root
	for i := 1; i <= 32; i++ {
		child := (key) >> 31 & 0x1
		key = key << 1
		//fmt.Printf("child %b key %032b\n", child, key)
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
	/*
		if root.Value != -1 {
			valueList = append(valueList, root.Value)
		}
	*/
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
