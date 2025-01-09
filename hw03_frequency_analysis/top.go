package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type NodeValue struct {
	Word  string
	Count int
}

type Node struct {
	Left  *Node
	Value NodeValue
	Right *Node
}

func Insert(node *Node, value NodeValue) {
	if value.Count > node.Value.Count {
		if node.Right == nil {
			node.Right = &Node{nil, value, nil}
		} else {
			Insert(node.Right, value)
		}
	}
	if value.Count <= node.Value.Count {
		if node.Left == nil {
			node.Left = &Node{nil, value, nil}
		} else {
			Insert(node.Left, value)
		}
	}
}

func MakeSortedArray(node *Node, nums *[]NodeValue) {
	if node.Right == nil {
		*nums = append(*nums, node.Value)
	} else {
		MakeSortedArray(node.Right, nums)
		*nums = append(*nums, node.Value)
	}

	if node.Left == nil {
		lIndex := len(*nums)
		if lIndex > 0 && node.Value.Word != ((*nums)[lIndex-1]).Word {
			*nums = append(*nums, node.Value)
		}
	} else {
		MakeSortedArray(node.Left, nums)
	}
}

func SortSubslices(values *[]NodeValue) []string {
	result := &[]string{}
	sl := &[]string{}
	var previous *NodeValue
	for i := range *values {
		nodeValue := (*values)[i]

		if previous == nil {
			previous = &nodeValue
			*sl = append(*sl, nodeValue.Word)
			continue
		}

		if previous.Count == nodeValue.Count {
			*sl = append(*sl, nodeValue.Word)
			previous = &nodeValue
		} else {
			sort.Strings(*sl)
			*result = append(*result, *sl...)
			previous = nil
			sl = &[]string{}
			*sl = append(*sl, nodeValue.Word)

			if len(*result) >= 10 {
				break
			}
		}
	}
	sort.Strings(*sl)
	*result = append(*result, *sl...)

	if len(*result) >= 10 {
		return (*result)[0:10]
	}
	return *result
}

func CleanAll(str string) string {
	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, "\t", " ")
	return str
}

func Top10(str string) []string {
	if len(str) == 0 {
		return []string{}
	}
	// "Чистим" строку
	str = CleanAll(str)
	// Считаем частоту слов
	wordsFrequency := make(map[string]int)
	for _, word := range strings.Split(str, " ") {
		if len(word) > 0 {
			wordsFrequency[word] += 1
		}
	}
	// Строим бинарное дерево
	var root *Node
	for k, v := range wordsFrequency {
		value := NodeValue{k, v}
		if root == nil {
			root = &Node{nil, value, nil}
			continue
		}
		Insert(root, value)
	}
	// Обход дерева в сортированный массив
	maxNums := &[]NodeValue{}
	MakeSortedArray(root, maxNums)
	// Сортируем слова с одинаковой частотой
	result := SortSubslices(maxNums)
	return result
}
