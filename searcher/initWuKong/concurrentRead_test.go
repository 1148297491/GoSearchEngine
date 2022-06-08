package initwukong

import (
	"GoSearchEngine/searcher"
	"GoSearchEngine/searcher/utils"
	"GoSearchEngine/searcher/words"
	"log"
	"testing"
)

func TestMultiRead(t *testing.T) {

	wordTokenizer := words.NewTokenizer("/home/tanlin/golang/src/GoSearchEngine/searcher/words/data/dictionary.txt")

	var engine = &searcher.Engine{
		IndexPath: "../tests/indexTest/test3", // 索引文件路径
		Tokenizer: wordTokenizer,              //定义分词器
	}

	//设置engine配置文件
	option := engine.GetOptions()

	// 引擎一直存在索引实体消费者协程
	engine.InitOption(option)

	filePath := "/home/tanlin/golang/src/GoSearchEngine/wukong.csv"
	read := NewReadInstance(engine.Shard)

	_time := utils.ExecTime(func() {
		read.MultiRead(filePath, engine, 5)

		for engine.GetQueue() > 0 {
			// 阻塞
		}
	})

	log.Printf("初始化耗时为: {%+v} ms", _time)
}

func TestSplitCsvLine(t *testing.T) {
	var str = "https://pic.rmb.bdstatic.com/19539b3b1a7e1daee93b0f3d99b8e795.png,\"曾是名不见经传的王平,为何能够取代魏延,成为蜀汉\""
	read := NewReadInstance(1)
	res := read.splitCsvLine([]byte(str))
	log.Printf("切分后长度为{%+v}, 切分答案为：%+v", len(res), res)

	var str2 = "\"https://gimg2.baidu.com/image_search/src=http%3A%2F%2F5b0988e595225.cdn.sohucs.com%2Fimages%2F20200326%2Fffc00cb6bc944e5b9ab2673c4873b24c.jpeg&refer=http%3A%2F%2F5b0988e595225.cdn.sohucs.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=jpeg?sec=1632531279&t=1b9ba84f70ddebdda6601a5576d37c50\",\"美沃可视数码裂隙灯,检查眼前节健康状况\""
	res = read.splitCsvLine([]byte(str2))

	log.Printf("切分后长度为{%+v}, 切分答案为：%+v", len(res), res)
}
