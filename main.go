package main

import (
	"GoSearchEngine/middleware"
	"GoSearchEngine/routers"
	"GoSearchEngine/searcher"
	"GoSearchEngine/searcher/model"
	"GoSearchEngine/searcher/words"
	"GoSearchEngine/web"
	"encoding/csv"

	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

func GetHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() * 2) //定义程序线程数

	middleware.InitConfig()
	middleware.InitDb()

	r := routers.SetupRouter()
	middleware.InitParamValidation()

	// 初始化引擎
	wordTokenizer := words.NewTokenizer("./searcher/words/data/dictionary.txt")

	var engine = &searcher.Engine{
		IndexPath: "./tests/indexTest/test0", // 索引文件路径
		Tokenizer: wordTokenizer,             //定义分词器
	}

	//设置engine配置文件
	option := engine.GetOptions()

	// 引擎一直存在索引实体消费者协程
	engine.InitOption(option)

	webSearch := &web.WebSearch{
		SearchEngine: engine,
	}

	initWukong(webSearch)

	counts := 1

	c := cron.New()
	spec := "0 */60 * * * ?" //设定60分钟检测一次
	err := c.AddFunc(spec, func() {
		if len(webSearch.SearchEngine.DeleteSet) > int(float32(webSearch.SearchEngine.GetDocumentCount())*0.1) {
			//设定检测规则：如果删除的文章数量大于当前数据库数量的百分之10，重新更换引擎数据库，并将其保存为一个新版本
			wg := &sync.WaitGroup{}
			wg.Add(1)
			counts++
			newEngine := webSearch.SearchEngine.TransformToNewEngine(counts, wg)
			wg.Wait()
			webSearch.SearchEngine = newEngine

			log.Println("重新初始化引擎！")
		}
	})
	if err != nil {
		log.Printf("AddFunc error : %v", err)
		return
	}
	c.Start()

	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": "引擎在线",
		})
	})
	r.POST("/search", webSearch.SearchProcessFun)

	err1 := r.Run(":9633") // 指定一个监听端口
	if err1 != nil {
		log.Printf("AddFunc error : %v", err)
		return
	}

	defer c.Stop()
}

func initWukong(webSearch *web.WebSearch) {
	openfile, err := os.Open("./wukong50k.csv")

	if err != nil {
		log.Printf("打开wukong50k文件失败, err=[%v]", err)
	}

	defer openfile.Close()

	readCsv := csv.NewReader(openfile)
	dataIndex := 0
	id := uint32(0)
	// 控制100个协程来生产索引
	producerChannel := make(chan struct{}, 10)
	defer close(producerChannel)

	for {
		csvLine, err := readCsv.Read()
		if err == io.EOF {
			// 读到文件结尾退出
			break
		}
		if dataIndex == 0 {
			dataIndex++
			continue
		}

		data := make(map[string]interface{})
		data["id"] = id
		data["url"] = csvLine[0]
		data["caption"] = csvLine[1]

		idxInstance := model.IndexDoc{
			Id:       id,
			Text:     csvLine[1],
			Document: data,
		}
		// log.Printf("向引擎的管道发送索引实体：%+v\n", idxInstance)
		// 创建生产者 放入索引实体 最多并发100个协程
		producerChannel <- struct{}{}
		go func() {
			webSearch.SearchEngine.IndexDocument(&idxInstance)
			<-producerChannel
		}()
		id++
	}

	// 到这一步消费者异步的消费不一定完成
	// 尝试获取 第一条数据的第一个切分词

	for webSearch.SearchEngine.GetInitQueue() > 0 {
		// 阻塞保证文档处理完成
		// log.Printf("索引队列长度：%+v", engine.GetQueue())
	}

	log.Println("初始化悟空数据集完成")
}
