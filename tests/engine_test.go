package tests

import (
	"encoding/csv"
	"fmt"
	"gofound/searcher"
	"gofound/searcher/model"
	"gofound/searcher/words"
	"io"
	"log"
	"os"
	"testing"
)

func TestEngineIndex(t *testing.T) {
	// 测试搜索引擎检索
	// 创建分词器
	tokenizer := words.NewTokenizer("../searcher/words/data/dictionary.txt")

	// 创建引擎
	var engine = &searcher.Engine{
		IndexPath: "../tests/indexTest/test2", // 索引文件路径
		Tokenizer: tokenizer,
	}

	option := engine.GetOptions()

	// 引擎一直存在索引实体消费者协程
	engine.InitOption(option)

	// 独取目标文件，创建索引实体
	// 独取cvs文件
	openfile, err := os.Open("../wukong50k.csv")

	if err != nil {
		log.Printf("打开wukong50k文件失败, err=[%v]", err)
	}

	defer openfile.Close()

	// 创建csv读取接口实现
	// csv文件第一行是标题，直接不处理
	readCsv := csv.NewReader(openfile)
	dataIndex := 0
	id := uint32(0)
	// 控制100个协程来生产索引
	producerChannel := make(chan struct{}, 100)
	defer close(producerChannel)

	for {
		// log.Println("!!")
		csvLine, err := readCsv.Read()
		if err == io.EOF {
			// 读到文件结尾退出
			break
		}
		if dataIndex == 0 {
			dataIndex++
			continue
		}
		// if id == 5000 {
		// 	break
		// }
		// fmt.Printf("%T", csvLine)
		id++
		data := make(map[string]interface{})
		data["id"] = id
		data["url"] = csvLine[0]
		data["caption"] = csvLine[1]

		// 创建索引实体
		/*
			// IndexDoc 索引实体
			type IndexDoc struct {
				Id       uint32                 `json:"id,omitempty"`
				Text     string                 `json:"text,omitempty"`
				Document map[string]interface{} `json:"document,omitempty"`
			}
		*/
		idxInstance := model.IndexDoc{
			Id:       id,
			Text:     csvLine[1],
			Document: data,
		}
		// log.Printf("向引擎的管道发送索引实体：%+v\n", idxInstance)
		// 创建生产者 放入索引实体 最多并发100个协程
		producerChannel <- struct{}{}
		go func() {
			engine.IndexDocument(&idxInstance)
			<-producerChannel
		}()
	}

	// 到这一步消费者异步的消费不一定完成
	// 尝试获取 第一条数据的第一个切分词
	for engine.GetQueue() > 0 {
		// 阻塞保证文档处理完成
		// log.Printf("索引队列长度：%+v", engine.GetQueue())
	}

	splitWords := tokenizer.Cut("美沃可视数码裂隙灯,检查眼前节健康状况")
	log.Printf("分词结果是：%+v", splitWords)

	for _, word := range splitWords {
		idxList, isFind := engine.GetInvertedIdxListByWord(word)
		if isFind == false {
			log.Printf("关键词[%s]找不到\n", word)
			continue
		}
		log.Printf("关键词[%s]的倒排索引文档列表是%+v\n", word, idxList)
	}

	searchReq := &model.SearchRequest{
		Query:       "",
		Order:       "",
		Page:        1,
		Limit:       10,
		FilterWords: []string{},
	}

	res := engine.MultiSearch(searchReq)
	log.Printf("结果是：%+v", res)
}

func TestReadCsv(t *testing.T) {
	openfile, err := os.Open("../wukong50k.csv")

	if err != nil {
		log.Printf("打开wukong50k文件失败, err=[%v]", err)
	}

	defer openfile.Close()

	// 创建csv读取接口实例
	readCsv := csv.NewReader(openfile)
	for i := 0; i < 2; i++ {
		csvLine, _ := readCsv.Read()
		fmt.Printf("%T\n", csvLine)
		fmt.Println(csvLine)
	}
}
