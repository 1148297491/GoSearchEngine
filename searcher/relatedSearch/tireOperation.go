package relatedsearch

import (
	"bytes"
	"gofound/searcher/utils"
	"gofound/searcher/words"
)

const (
	maxLayers          = 10 //表示最多bfs多少层
	maxRelatedWordsLen = 10 //表示最多放几个related推荐
)

func initTokenizer(dictionaryPath string) *words.Tokenizer {
	return words.NewTokenizer(dictionaryPath)
}

func addString(str1 string, str2 string) string { //拼接字符串
	var buffer bytes.Buffer
	buffer.WriteString(str1)
	buffer.WriteString(str2)
	return buffer.String()
}
func JoinWords(str string, trie *Trie, tokenizer *words.Tokenizer) {
	//func JoinWords(str string, trie *Trie) {
	trie.AddConsumeChannel(str) //将关键词列表存入前缀树
}

func GetRelatedWords(str string, trie *Trie, tokenizer *words.Tokenizer) []string {
	inputString := tokenizer.Cut(str, true) //使用精准分词
	existString := make(map[uint32]bool)
	q := new(Queue) //初始化队列
	q.Init()

	var tempWordsList []string
	var tempWords string
	for _, word := range inputString {
		tempWordsList = append(tempWordsList, word) //每次多加一个关键字放入队列中
		tempWords = addString(tempWords, word)
		existString[utils.StringToInt(tempWords)] = true

		q.Enqueue(NewQueueNode(trie.FindWordNode(tempWordsList), tempWords, maxLayers))
	}

	var relatedWords []string //推荐得到的相关搜索内容

	for {
		if q.Size() == 0 || len(relatedWords) >= maxRelatedWordsLen {
			//如果队列为空（表示推荐已结束）或者推荐超过了指定长度则退出
			break
		}
		frontNode := q.Dequeue().(*QueueNode)               //获取队列头节点
		if frontNode.Layers == 0 || frontNode.Node == nil { //达到最大bfs深度或无叶子结点则跳过
			continue
		}

		for _, childrenTrieNode := range frontNode.Node.Children { //遍历当前结点的所有孩子结点，进行bfs
			nextStrings := frontNode.Strings //深拷贝当前的bfs字符串，将其增加该孩子结点的内容，遍历层数-1放入队列中
			nextStrings = addString(nextStrings, childrenTrieNode.NodeString)
			_, ok := existString[utils.StringToInt(nextStrings)]
			if ok {
				// 如果是用户输入搜索的前缀则直接跳过
				continue
			}
			if childrenTrieNode.IsEnding {
				relatedWords = append(relatedWords, nextStrings) //如果是结尾词则该推荐有效，可以放入推荐列表中
			}
			if len(relatedWords) >= maxRelatedWordsLen {
				break
			}
			q.Enqueue(NewQueueNode(childrenTrieNode, nextStrings, frontNode.Layers-1)) //将修改后的结点入队
		}
	}
	for { //清空内存
		if q.Size() != 0 {
			q.Dequeue()
		} else {
			break
		}
	}
	return relatedWords
}
