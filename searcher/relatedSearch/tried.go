package relatedsearch

import (
	"fmt"
	"gofound/searcher/utils"
	"gofound/searcher/words"
)

// Trie 树节点
type TrieNode struct {
	NodeString string // Unicode 字符
	IsEnding   bool   // 是否是单词结尾
	CurrentNum uint16
	Children   map[uint32]*TrieNode // 该节点的子节点
}

// 初始化 Trie 树节点
func NewTrieNode(word string) *TrieNode {
	return &TrieNode{
		NodeString: word,
		IsEnding:   false,
		CurrentNum: 0,
		Children:   make(map[uint32]*TrieNode),
	}
}

// Trie 树结构
type Trie struct {
	root            *TrieNode // 根节点指针
	trieNodeChannel chan string
	trieTokenizer   *words.Tokenizer
}

// 初始化 Trie 树
func NewTrie(tokenizer *words.Tokenizer) *Trie {
	// 初始化根节点
	Node := NewTrieNode("0")
	tire := &Trie{
		root:            Node,
		trieNodeChannel: make(chan string, 1000),
		trieTokenizer:   tokenizer,
	}

	go func() {
		tire.ConsumeChannel()
	}()

	return tire
}

func (t *Trie) ConsumeChannel() {
	for text := range t.trieNodeChannel {
		t.Insert(text)
	}
}

// 往 Trie 树中插入一个多关键字句子
func (t *Trie) Insert(text string) {
	words := t.trieTokenizer.Cut(text, []string{"0"})
	node := t.root               // 获取根节点
	for _, word := range words { // 将该句子中的每个关键词以此放入前缀树
		code := utils.StringToInt(word)
		value, ok := node.Children[code] // 获取 code 编码对应子节点
		if !ok {
			// 不存在则初始化该节点
			value = NewTrieNode(word)
			// 然后将其添加到子节点字典
			node.Children[code] = value
		}
		node.Children[code].CurrentNum++
		// 当前节点指针指向当前子节点
		node = value
	}
	node.IsEnding = true // 一个单词遍历完所有字符后将结尾字符打上标记
}

// 在 Trie 树中查找一个单词
func (t *Trie) Find(words []string) bool {
	node := t.root
	for _, word := range words {
		code := utils.StringToInt(word)
		value, ok := node.Children[code] // 获取 code 编码对应子节点
		if !ok {
			// 不存在则直接返回
			return false
		}
		fmt.Println(word, node.Children[code].CurrentNum)
		// 否则继续往后遍历
		node = value
	}
	if node.IsEnding == false {
		return false // 不能完全匹配，只是前缀
	}
	return true // 找到对应单词
}

//在trie树中找到某个单词组最后的那个结点
func (t *Trie) FindWordNode(words []string) *TrieNode {
	node := t.root
	for _, word := range words {
		code := utils.StringToInt(word)
		value, ok := node.Children[code] // 获取 code 编码对应子节点
		if !ok {
			//找不到该结点就是nil
			return nil
		}
		node = value
	}
	return node //找到对应的最后一个词的结点
}

func (t *Trie) AddConsumeChannel(words string) {
	t.trieNodeChannel <- words
	/*for _, word := range words {
		t.trieNodeChannel <- word
	}*/
}
