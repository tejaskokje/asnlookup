package asnlookup

import (
	"fmt"
)

type Node struct {
	Label uint32
	// value is asn for this prefix
	Value int
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
		Value: -1,
		Left:  nil,
		Right: nil,
	}

}
func Insert(t *Trie, key uint32, value int, prefixLen int) {
	fmt.Println("Inserting ", key, value, prefixLen)
	root := t.Root
	for i := 0; i < prefixLen; i++ {
		child := (key << uint32(i)) & 1
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

	root.Label = key
	root.Value = value

	return
}

func Find(t *Trie, key uint32) []int {

	valueList := []int{}
	root := t.Root
	for i := 0; i < 32; i++ {
		child := (key << uint32(i)) & 1
		if child == 0 && root.Left != nil {
			if root.Value != -1 {
				valueList = append(valueList, root.Value)
			}
			root = root.Left
		} else if child == 1 && root.Right != nil {

			if root.Value != -1 {
				valueList = append(valueList, root.Value)
			}

			root = root.Right
		} else {
			break
		}
	}

	if root.Value != -1 {
		valueList = append(valueList, root.Value)
	}

	return valueList
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

		fmt.Println(relation, n.Label, n.Value)
		DumpNode(n.Left, 0)
		DumpNode(n.Right, 1)
	}
}
