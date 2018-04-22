package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var huffmanCodeTable map[uint16]string

func init() {

	huffmanCodeTable = make(map[uint16]string)
	rand.Seed(time.Now().UnixNano())
}

type Tree interface {
	Freq() int
}

type Leaf struct {
	frequency int
	value     uint16
	code      string
}

type Node struct {
	frequency   int
	left, right Tree
}

func (leaf Leaf) Freq() int {
	return leaf.frequency
}

func (node Node) Freq() int {
	return node.frequency
}

type impTree []Tree

func (th impTree) Len() int {
	return len(th)
}
func (th impTree) Less(i, j int) bool {
	return th[i].Freq() < th[j].Freq()
}

func (th impTree) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
}
func (th *impTree) Push(ele interface{}) {
	*th = append(*th, ele.(Tree))
}

func (th *impTree) Pop() (popped interface{}) {
	popped = (*th)[len(*th)-1]
	*th = (*th)[:len(*th)-1]
	return
}

func buildTree(symbolfrequencyuency map[uint16]int) Tree {

	var trees impTree
	for c, f := range symbolfrequencyuency {
		trees = append(trees, Leaf{f, c, ""})
	}
	heap.Init(&trees)
	for trees.Len() > 1 {
		// two trees with least frequencyuency
		a := heap.Pop(&trees).(Tree)
		b := heap.Pop(&trees).(Tree)

		// put into new node and re-insert into queue
		heap.Push(&trees, Node{a.Freq() + b.Freq(), a, b})
	}

	return heap.Pop(&trees).(Tree)
}

func buildCodes(tree Tree, prefix []byte) {

	switch i := tree.(type) {
	case Leaf:
		i.code = string(prefix)
		huffmanCodeTable[i.value] = i.code
	case Node:

		// build left
		prefix = append(prefix, '0')
		buildCodes(i.left, prefix)
		prefix = prefix[:len(prefix)-1]

		// build right
		prefix = append(prefix, '1')
		buildCodes(i.right, prefix)
		prefix = prefix[:len(prefix)-1]
	}

}
func Benchmarkfrequency(b *testing.B) {

	symbolfrequencyuency := make(map[uint16]int)

	for i := 0; i < 10; i++ {
		symbolfrequencyuency[uint16(i)] = rand.Intn(30)

	}
	fmt.Println(symbolfrequencyuency)
	tree := buildTree(symbolfrequencyuency)
	fmt.Println("SYMBOL\tWEIGHT\tHUFFMAN CODE")
	buildCodes(tree, []byte{})
	fmt.Println(huffmanCodeTable)
}
