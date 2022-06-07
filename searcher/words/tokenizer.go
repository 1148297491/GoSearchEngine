package words

import (
	"embed"
	"GoSearchEngine/searcher/utils"
	"log"
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

// Option可选参数, 传入有且仅一个时候判断类型, bool类型做进相关搜索的精细分词, slice表示搜索语句过滤
func (t *Tokenizer) Cut(text string, option ...interface{}) []string {
	//不区分大小写
	text = strings.ToLower(text)
	//移除所有的标点符号
	text = utils.RemovePunctuation(text)
	//移除所有的空格
	text = utils.RemoveSpace(text)

	for _, op := range option {
		switch optype := op.(type) {
		case bool:
			return t.preciseCutforRelatedSearch(text)
		case []string:
			return t.cutWithfilterWords(text, op.([]string))
		default:
			log.Printf("[Tokenizer.Cut]可选参数传递错误类型,具体为: %+v", optype)
		}

		break // 保证仅处理第一个参数
	}

	// 没有可选参数 或者可选参数错误 按照没有可选参数进行处理
	var wordSet = make(utils.Set)
	var wordsSlice []string
	var resultChan = t.seg.CutForSearch(text, true)
	var stopWordsSet = GetStopWordsSet()

	//对用户搜索的内容进行相应的处理（屏蔽搜索词以及去重）
	for {
		word, ok := <-resultChan
		if !ok {
			break
		}

		if stopWordsSet.IsExist(word) {
			continue
		}

		if wordSet.IsExist(word) {
			continue
		}

		// 未存在重复的词语放进返回结果列表
		wordsSlice = append(wordsSlice, word)
	}

	return wordsSlice
}

// 当采用精准分词模式返回其分词结果 相关搜索使用
func (t *Tokenizer) preciseCutforRelatedSearch(text string) []string {
	var wordsSlice []string
	var resultChan = t.seg.Cut(text, true)
	var stopWordsSet = GetStopWordsSet()

	for {
		w, ok := <-resultChan

		if !ok {
			break
		}

		if stopWordsSet.IsExist(w) {
			continue
		}

		wordsSlice = append(wordsSlice, w)
	}
	return wordsSlice
}

// 带有过滤词的分词方式
func (t *Tokenizer) cutWithfilterWords(text string, filterWords []string) []string {
	var wordsSlice []string
	var resultChan <-chan string
	var wordSet = make(utils.Set)
	var filterWordsMap = make(map[string]int, len(filterWords))
	var stopWordsSet = GetStopWordsSet()

	// 先做一个位置标记
	for idx, word := range filterWords {
		filterWordsMap[word] = idx
	}

	resultChan = t.seg.CutForSearch(text, true)

	for {
		word, ok := <-resultChan
		if !ok {
			break
		}

		// 考虑到搜索引擎分词可能初夏

		// 过滤词中存在
		if _, ok := filterWordsMap[word]; ok {
			continue
		}

		if stopWordsSet.IsExist(word) {
			continue
		}

		if wordSet.IsExist(word) {
			continue
		}

		// 未存在重复的词语放进返回结果列表
		wordsSlice = append(wordsSlice, word)
	}
	return wordsSlice
}
