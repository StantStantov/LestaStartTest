package huffman

import (
	"container/heap"
	"fmt"
	"io"
	"strings"
)

func Encode(data io.Reader) (string, map[string]string, error) {
	buffer, err := io.ReadAll(data)
	if err != nil {
		return "", nil, fmt.Errorf("huffman/huffman.Encode: [%w]", err)
	}

	codesTable := createCodesTable(createHuffmanTree(createFrequencyTable(buffer)))

	builder := new(strings.Builder)
	for _, rune := range buffer {
		code, ok := codesTable[string(rune)]
		if !ok {
			return "", nil, fmt.Errorf("huffman/huffman.Encode: [%c can't be properly encoded]", rune)
		}
		builder.WriteString(code)
	}

	return builder.String(), codesTable, nil
}

func createFrequencyTable(data []byte) map[byte]uint {
	frequencies := map[byte]uint{}
	for _, rune := range data {
		frequencies[rune]++
	}

	return frequencies
}

func createHuffmanTree(frequencyTable map[byte]uint) *huffmanNode {
	weights := &huffmanHeap{}
	heap.Init(weights)
	for value, frequency := range frequencyTable {
		heap.Push(weights, &huffmanNode{Value: value, Frequency: frequency})
	}

	for weights.Len() > 1 {
		left := heap.Pop(weights).(*huffmanNode)
		right := heap.Pop(weights).(*huffmanNode)

		heap.Push(weights, &huffmanNode{Left: left, Right: right, Frequency: left.Frequency + right.Frequency})
	}

	root := heap.Pop(weights).(*huffmanNode)
	return root
}

func createCodesTable(huffmanTree *huffmanNode) map[string]string {
	codesTable := make(map[string]string)
	if huffmanTree == nil {
		return codesTable
	}
	generateCode(huffmanTree, "", codesTable)

	return codesTable
}

func generateCode(node *huffmanNode, prefix string, codesTable map[string]string) {
	if node.Left == nil && node.Right == nil {
		codesTable[string(node.Value)] = prefix
		return
	}
	if node.Left != nil {
		generateCode(node.Left, prefix+"0", codesTable)
	}
	if node.Right != nil {
		generateCode(node.Right, prefix+"1", codesTable)
	}
}

type huffmanHeap []*huffmanNode

func (h huffmanHeap) Len() int {
	return len(h)
}

func (h huffmanHeap) Less(a, b int) bool {
	return h[a].Frequency < h[b].Frequency
}

func (h huffmanHeap) Swap(a, b int) {
	h[a], h[b] = h[b], h[a]
}

func (h *huffmanHeap) Push(a any) {
	*h = append(*h, a.(*huffmanNode))
}

func (h *huffmanHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type huffmanNode struct {
	Parent    *huffmanNode
	Left      *huffmanNode
	Right     *huffmanNode
	Value     byte
	Frequency uint
}
