package words

import (
	"log"
	"testing"
)

func TestCut(t *testing.T) {
	wordTokenizer := NewTokenizer("/home/tanlin/golang/src/GoSearchEngine/searcher/words/data/dictionary.txt")
	// fileterWord := []string{"谭琳", "哈尔滨"}

	res := wordTokenizer.Cut("哈尔滨工程大学", true)

	log.Printf("分词的结果是：%+v", res)
	// log.Printf("处理后的过滤词是：%+v", fileterWord)
}
