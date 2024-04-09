package fuzzysearch

import "diggle/searcher/models"

type Node struct {
	item     models.WordFrequency
	children map[int]*Node
}

type BKTree struct {
	root *Node
}

type Result struct {
	word      string
	distance  int
	frequency int
}

func NewBKTree() *BKTree {
	return &BKTree{root: nil}
}

func (b *BKTree) Add(wordFreq models.WordFrequency) bool {

	if b.root == nil {
		b.root = &Node{wordFreq, map[int]*Node{}}
		return true
	}

	currentNode := b.root

	for {
		distance := GetEditDistance(wordFreq.Word, currentNode.item.Word)
		if distance == 0 {
			return false
		}

		if node, ok := currentNode.children[int(distance)]; ok {
			currentNode = node
			continue
		}

		node := &Node{wordFreq, map[int]*Node{}}
		currentNode.children[int(distance)] = node
		break
	}

	return true
}

func (b *BKTree) Search(wordFreq string, tolerance int) []Result {

	result := []Result{}
	if b.root == nil {
		return result
	}

	var traverse func(node *Node)

	traverse = func(node *Node) {
		distance := GetEditDistance(wordFreq, node.item.Word)
		if distance == 0 {
			result = append(result, Result{node.item.Word, 0, node.item.Count})
			return
		}
		if distance <= tolerance {
			result = append(result, Result{node.item.Word, distance, node.item.Count})
		}
		for i := distance - tolerance; i <= distance+tolerance; i++ {
			if child, ok := node.children[i]; ok {
				traverse(child)
			}
		}
	}

	traverse(b.root)

	return result
}
