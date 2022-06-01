package words

import (
	"embed"
	"gofound/searcher/utils"
	"strings"

	"github.com/wangbin/jiebago"
)

var (
	//go:embed data/*.txt
	dictionaryFS embed.FS
)

type Tokenizer struct {
	seg jiebago.Segmenter
}

func NewTokenizer(dictionaryPath string) *Tokenizer {
	//打开词典列表构建分词器
	file, err := dictionaryFS.Open("data/dictionary.txt")
	if err != nil {
		panic(err)
	}
	utils.ReleaseAssets(file, dictionaryPath)

	tokenizer := &Tokenizer{}

	err = tokenizer.seg.LoadDictionary(dictionaryPath)
	if err != nil {
		panic(err)
	}

	return tokenizer
}

func (t *Tokenizer) Cut(text string, filterWords ...[]string) []string {
	//不区分大小写
	text = strings.ToLower(text)
	//移除所有的标点符号
	text = utils.RemovePunctuation(text)
	//移除所有的空格
	text = utils.RemoveSpace(text)

	var wordSet = make(utils.Set)
	var filterWordsMap = make(map[string]int)
	var wordsSlice []string
	var resultChan <-chan string

	if len(filterWords) > 0 && len(filterWords[0]) == 1 && filterWords[0][0] == "0" {
		//当采用精准分词模式返回其分词结果
		resultChan = t.seg.Cut(text, true)
		for {
			word, ok := <-resultChan
			if !ok {
				break
			}
			if _, ok := stopWordsMap[word]; !ok {
				wordsSlice = append(wordsSlice, word)
			}
		}
		return wordsSlice
	}

	//当存在屏蔽词时，首先进行标记
	if len(filterWords) > 0 {
		for idx, word := range filterWords[0] {
			filterWordsMap[word] = idx
		}
	}

	resultChan = t.seg.CutForSearch(text, true)

	//对用户搜索的内容进行相应的处理（屏蔽搜索词以及去重）
	for {
		word, ok := <-resultChan
		if !ok {
			break
		}

		// 对重复词语 和 过滤词语进行筛除
		if _, ok := filterWordsMap[word]; ok && len(filterWords) > 0 {
			filterWords[0][filterWordsMap[word]] = "" // 空的时候不对
			continue
		}

		if wordSet.IsExist(word) {
			continue
		}

		// 未存在或者重复的词语放进返回结果列表
		if _, ok := stopWordsMap[word]; !ok {
			wordsSlice = append(wordsSlice, word)
			/*if word == "是" {
				fmt.Println(stopWordsMap["是"])
			}*/
		}
	}

	return wordsSlice
}
