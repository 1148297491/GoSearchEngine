package searcher

import (
	"fmt"
	"gofound/searcher/arrays"
	"gofound/searcher/model"
	"gofound/searcher/pagination"
	relatedsearch "gofound/searcher/relatedSearch"
	"gofound/searcher/sorts"
	"gofound/searcher/storage"
	"gofound/searcher/utils"
	"gofound/searcher/words"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	IndexPath string  //索引文件存储目录
	Option    *Option //配置 即方便初始化levelDB

	TrieRoot *relatedsearch.Trie // 相关搜索前缀树

	invertedIndexStorages []*storage.LeveldbStorage //倒排索引, key=word, value=[]indexs
	positiveIndexStorages []*storage.LeveldbStorage //正排索引, 一个文档id对应多个words
	docStorages           []*storage.LeveldbStorage //文档仓

	sync.Mutex //锁 互斥锁
	//initWG                *sync.WaitGroup        //LevelDB初始化信号量机制
	//IndexWG               *sync.WaitGroup        //索引处理信号量机制
	sync.WaitGroup                               //LevelDB初始化信号量机制
	addDocumentWorkerChan []chan *model.IndexDoc //添加索引的通道
	IsDebug               bool                   //是否调试模式
	Tokenizer             *words.Tokenizer       //分词器
	DatabaseName          string                 //数据库名

	Shard int //分片数
}

type Option struct {
	InvertedIndexName string //倒排索引
	PositiveIndexName string //正排索引
	DocIndexName      string //文档存储
}

const (
	PreTag  = "<mark>" //定义返回结果高亮的前缀
	PostTag = "</mark>" //定义返回结果高亮的后缀
)

// Init 初始化索引引擎（对一个具体的文档进行初始化）
func (e *Engine) Init() {
	// 调用waitGroup{} 实现信号量控制, 并发处理文档需要加锁
	e.Add(1)
	defer e.Done()

	e.TrieRoot = relatedsearch.NewTrie(e.Tokenizer)

	if e.Option == nil {
		e.Option = e.GetOptions()
	}

	log.Println("数据存储目录：", e.IndexPath) // 根据数据库相关文件存储位置

	e.addDocumentWorkerChan = make([]chan *model.IndexDoc, e.Shard) // 文件指定分片数，没有指定则默认10
	//初始化文件存储
	for shard := 0; shard < e.Shard; shard++ {

		//初始化chan
		worker := make(chan *model.IndexDoc, 1000)
		e.addDocumentWorkerChan[shard] = worker //定义索引缓冲区
		go e.DocumentWorkerExec(worker)         // 启动协程并发消费索引实体

		// levelDB按照文件名字+分片值
		s, err := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.DocIndexName, shard)))
		if err != nil {
			panic(err)
		}
		//存放文档相关的levelDB，levelDB存放文档实体：id+文档内容, 后面需要根据分片号获得对应的文档仓
		e.docStorages = append(e.docStorages, s)

		//初始化倒排索引数据库
		ks, kerr := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.InvertedIndexName, shard)))
		if kerr != nil {
			panic(err)
		}
		e.invertedIndexStorages = append(e.invertedIndexStorages, ks)

		//id和keys映射 正排索引
		iks, ikerr := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.PositiveIndexName, shard)))
		if ikerr != nil {
			panic(ikerr)
		}
		e.positiveIndexStorages = append(e.positiveIndexStorages, iks)
	}
	go e.automaticGC()
	log.Println("初始化完成")
}

// 自动保存索引，10秒钟检测一次
func (e *Engine) automaticGC() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		<-ticker.C
		//定时GC
		runtime.GC()
	}
}

func (e *Engine) IndexDocument(doc *model.IndexDoc) {
	//首先获得文档id对应的分片号, 然后向引擎相应管道发送索引实体
	// log.Printf("向引擎的管道发送索引实体：%+v", doc)
	e.addDocumentWorkerChan[e.getShard(doc.Id)] <- doc
}

// GetQueue 获取队列剩余
func (e *Engine) GetQueue() int {
	total := 0
	for _, v := range e.addDocumentWorkerChan {
		total += len(v)
	}
	return total
}

// DocumentWorkerExec 添加文档队列
func (e *Engine) DocumentWorkerExec(worker chan *model.IndexDoc) {
	// 容量1000的管道
	// 管道关闭 循环结束 消费管道
	for {
		doc := <-worker
		// log.Printf("管道获取索引实体：%+v\n", doc)
		e.AddDocument(doc)
	}
}

// getShard 计算索引分布在哪个文件块
func (e *Engine) getShard(id uint32) int {
	return int(id % uint32(e.Shard))
}

func (e *Engine) getShardByWord(word string) int {

	return int(utils.StringToInt(word) % uint32(e.Shard))
}

func (e *Engine) InitOption(option *Option) {

	if option == nil {
		//默认值
		option = e.GetOptions()
	}
	e.Option = option
	//shard默认值
	if e.Shard <= 0 {
		e.Shard = 10
	}
	//初始化其他的
	e.Init() //

}

func (e *Engine) getFilePath(fileName string) string {
	return e.IndexPath + string(os.PathSeparator) + fileName
}

func (e *Engine) GetOptions() *Option {
	return &Option{
		DocIndexName:      "docs",
		InvertedIndexName: "inverted_index",
		PositiveIndexName: "positive_index",
	}
}

func (e *Engine) AddDocument(index *model.IndexDoc) {
	//添加倒排索引+正排索引+文档仓库
	//等待初始化完成 引擎有不关闭的索引管道
	e.Wait()
	text := index.Text

	//采用引擎模式分词
	splitWords := e.Tokenizer.Cut(text)

	//添加该文章到前缀树中
	relatedsearch.JoinWords(text, e.TrieRoot, e.Tokenizer)

	//判断ID是否存在，如果存在，需要计算两次的差值，然后更新
	id := index.Id
	isUpdate := e.optimizeIndex(id, splitWords)

	//没有更新
	if !isUpdate {
		return
	}

	for _, word := range splitWords {
		e.addInvertedIndex(word, id)
	}

	//添加id索引
	e.addPositiveIndex(index, splitWords)
}

// 添加倒排索引
func (e *Engine) addInvertedIndex(word string, id uint32) {
	e.Lock()
	defer e.Unlock()

	shard := e.getShardByWord(word)

	s := e.invertedIndexStorages[shard]

	//string作为key
	key := []byte(word)

	//存在
	//添加到列表
	buf, find := s.Get(key)
	ids := make([]uint32, 0)
	if find {
		utils.Decoder(buf, &ids)
	}

	// id自增？
	if !arrays.BinarySearch(ids, id) {
		ids = append(ids, id)
	}

	s.Set(key, utils.Encoder(ids))
}

//	移除没有的词
func (e *Engine) optimizeIndex(id uint32, newWords []string) bool {
	//判断id是否存在
	e.Lock()
	defer e.Unlock()

	//计算差值
	removes, found := e.getDifference(id, newWords)
	if found && len(removes) > 0 {
		//从这些词中移除当前ID
		for _, word := range removes {
			e.removeIdInWordIndex(id, word)
		}
	}

	// 有没有更新
	return !found || len(removes) > 0

}

func (e *Engine) removeIdInWordIndex(id uint32, word string) {

	shard := e.getShardByWord(word)
	wordStorage := e.invertedIndexStorages[shard]

	//string作为key
	key := []byte(word)

	buf, found := wordStorage.Get(key)
	if found {
		ids := make([]uint32, 0)
		utils.Decoder(buf, &ids)

		//移除
		index := arrays.Find(ids, id)
		if index != -1 {
			ids = utils.DeleteArray(ids, index)
			if len(ids) == 0 {
				err := wordStorage.Delete(key)
				if err != nil {
					panic(err)
				}
			} else {
				wordStorage.Set(key, utils.Encoder(ids))
			}
		}
	}

}

// 计算差值
func (e *Engine) getDifference(id uint32, newWords []string) ([]string, bool) {

	shard := e.getShard(id)
	wordStorage := e.positiveIndexStorages[shard]
	key := utils.Uint32ToBytes(id)
	buf, found := wordStorage.Get(key)
	// 如果文档ID存在,取文档ID对应内容：id对应 一个[]string
	if found {
		oldWords := make([]string, 0)
		utils.Decoder(buf, &oldWords)

		//计算需要移除的
		removes := make([]string, 0)
		for _, word := range oldWords {

			//旧的在新的里面不存在，就是需要移除的
			if !arrays.ArrayStringExists(newWords, word) {
				removes = append(removes, word)
			}
		}
		return removes, true
	}

	return nil, false
}

// 添加正排索引,以及文档数据 id=>keys id=>doc
func (e *Engine) addPositiveIndex(index *model.IndexDoc, keys []string) {
	e.Lock()
	defer e.Unlock()

	indexByte := utils.Uint32ToBytes(index.Id)
	shard := e.getShard(index.Id)
	docStorage := e.docStorages[shard]

	//添加正排索引
	positiveIndexStorage := e.positiveIndexStorages[shard]

	doc := &model.StorageIndexDoc{
		IndexDoc: index,
		Keys:     keys,
	}

	//存储index及文章内容到文档仓库
	docStorage.Set(indexByte, utils.Encoder(doc))

	//设置index及分词列表到正排索引中
	positiveIndexStorage.Set(indexByte, utils.Encoder(keys))
}

// MultiSearch 多线程搜索
func (e *Engine) MultiSearch(request *model.SearchRequest) *model.SearchResult {
	//等待搜索初始化完成，即等待分片levelDB引擎初始化的完成
	e.Wait()
	// 分词 分词结果考虑需要过滤的词语
	// 首先需要对过滤词进行去重
	filterWords := e.DeleteDuplicatedWordAndCut(request.FilterWords)

	// query := e.deleteFilterWordsFromQuery(filterWords, request.Query)

	words := e.Tokenizer.Cut(request.Query, filterWords) // 切分搜索语句

	totalTime := float64(0)

	// 这个地方是对返回的结果的排序方案做出处理
	fastSort := &sorts.FastSort{
		IsDebug:      e.IsDebug,
		Order:        request.Order,
		FilterIdxSet: make(utils.Set),
		DataChannel:  make(chan uint32, 1000),
	}

	//并发消费倒排索引结果, 消费者（在并发情况下， 倒排索引放入管道，单协程消费，多协程生产）
	appendInvertedIndexWg := &sync.WaitGroup{}
	appendInvertedIndexWg.Add(1)
	go fastSort.AppendData(appendInvertedIndexWg)

	_time := utils.ExecTime(func() {
		// 根据分词结果获取倒排索引 结果装入fastSort中的管道
		e.searchGetInvertedIndex(fastSort, appendInvertedIndexWg, words, filterWords)
	})

	if e.IsDebug {
		log.Println("数组查找耗时：", totalTime, "ms")
		log.Println("搜索时间:", _time, "ms")
	}
	// 检查请求中的页码、每页限制、排序方式 如果为空采用默认值
	request = request.GetAndSetDefault()

	//计算交集得分和去重 倒排索引出现次数越多 说明包含搜索关键词越多 得分越高
	fastSort.Process()

	// 高亮词处理, 即返回结果中, 对包含的搜索关键词做出处理
	hightWords := make(utils.Set)

	for _, word := range words {
		hightWords.Add(word)
	}

	var result = &model.SearchResult{
		Total:         fastSort.Count(),                                                      // 结果集的长度大小
		Page:          request.Page,                                                          // 返回的结果页数
		Limit:         request.Limit,                                                         // 每页的限制大小
		RelatedSearch: relatedsearch.GetRelatedWords(request.Query, e.TrieRoot, e.Tokenizer), //获得相关搜索内容
	}

	_time += utils.ExecTime(func() {
		// 根据倒排索引获取目标数据
		e.getFinalSearchRes(fastSort, request, result, hightWords)
	})

	if e.IsDebug {
		log.Println("处理数据耗时：", _time, "ms")
	}
	result.Time = _time

	return result
}

func (e *Engine) searchGetInvertedIndex(fastSort *sorts.FastSort, appendInvertedIndexWg *sync.WaitGroup, words, filterWords []string) {
	base := len(words)
	wg := &sync.WaitGroup{}      //并发添加倒排索引, 生产者
	filterWg := sync.WaitGroup{} // 并发搜索过滤词的倒排索引，放入set中, 相当于消费者
	filterWg.Add(1)
	// 	按照切分的关键词进行检索
	// 同时需要考虑现有关键词的文档列表 可能包含用户过滤词 需要进一步筛查
	var filterWordIdxChannel = make(chan uint32, 100)
	go func() {
		defer filterWg.Done()
		for inf := range filterWordIdxChannel {
			fastSort.FilterIdxSet.Add(inf) //并发修改管道, 添加过滤词的文档index至set中, = 消费者
		}
	}()
	// 如果有需要过滤的词语需要进行处理
	for _, filterWord := range filterWords {
		if filterWord == "" {
			continue
		}
		wg.Add(1)
		go e.produceFilterIdx(filterWordIdxChannel, filterWord, wg) //生产屏蔽词index至set中
	}
	wg.Wait()
	// 生产者生产结束 关闭管道 等待消费结束
	close(filterWordIdxChannel)
	filterWg.Wait()

	wg.Add(base)
	for _, word := range words {
		// 查找的并发append 需要考虑并发安全问题
		go e.processKeySearch(word, fastSort, wg, base)
	}
	wg.Wait()
	// 生产完了所有的id
	close(fastSort.DataChannel)
	appendInvertedIndexWg.Wait()
}

func (e *Engine) getFinalSearchRes(fastSort *sorts.FastSort, request *model.SearchRequest, result *model.SearchResult, hightWords utils.Set) {
	pager := new(pagination.Pagination)

	pager.Init(request.Limit, fastSort.Count()) // 根据查询出来的文档数 获取返回结果的页数 count / limit

	//设置总页数
	result.PageCount = pager.PageCount

	// 对请求页码进行判断，判断是否请求内容是否在允许页码内，不在则直接返回
	if result.Page > pager.PageCount || request.Page < 0 {
		result.Documents = make([]model.ResponseDoc, 1)
		result.Documents[0] = model.ResponseDoc{
			Score:         -1,
			HighlightText: "<mark>没有更多内容了</mark>",
			OriginalText:  "没有更多内容了",
		}
		return
	}

	//读取单页的id
	if pager.PageCount != 0 {

		start, end := pager.GetPage(request.Page) //请求中的页码

		var resultItems = make([]model.SliceItem, 0)
		fastSort.GetAll(&resultItems, start, end)

		count := len(resultItems)

		result.Documents = make([]model.ResponseDoc, count)
		//只读取前面100个
		wg := new(sync.WaitGroup)
		wg.Add(count)
		for index, item := range resultItems {
			go e.getResult(item, &result.Documents[index], request, hightWords, wg)
		}
		wg.Wait()
	}
}

func (e *Engine) getResult(item model.SliceItem, doc *model.ResponseDoc, request *model.SearchRequest, hightWords utils.Set, wg *sync.WaitGroup) {
	buf := e.GetDocById(item.Id)
	defer wg.Done()
	doc.Score = item.Score

	if buf != nil {
		//gob解析
		storageDoc := new(model.StorageIndexDoc)
		utils.Decoder(buf, &storageDoc)
		doc.Url = storageDoc.Document["url"].(string) // 返回的url内容
		text := storageDoc.Text

		//处理关键词高亮 就是文档中包含的关键字做一个标记
		for _, key := range storageDoc.Keys {
			if hightWords.IsExist(key) {
				text = strings.ReplaceAll(text, key, fmt.Sprintf("%s%s%s", PreTag, key, PostTag))
			}
		}
		//放置原始文本
		doc.OriginalText = storageDoc.Text
		doc.HighlightText = text
	}

}

func (e *Engine) processKeySearch(word string, fastSort *sorts.FastSort, wg *sync.WaitGroup, base int) {
	defer wg.Done()

	shard := e.getShardByWord(word)
	//读取id：首先获取索引相关的levelDB,然后获取id列表
	invertedIndexStorage := e.invertedIndexStorages[shard]
	key := []byte(word)

	buf, find := invertedIndexStorage.Get(key)
	if find {
		ids := make([]uint32, 0)
		//解码
		utils.Decoder(buf, &ids)
		fastSort.Add(&ids) // 管道
	}
}

// GetInvertedIdxListByWord 仅获取倒排索引
func (e *Engine) GetInvertedIdxListByWord(word string) ([]uint32, bool) {
	// 获取levelDB分片号码
	shard := int(utils.StringToInt(word) % uint32(e.Shard))

	// 获取存储
	invertedIdxDB := e.invertedIndexStorages[shard]

	key := []byte(word)

	//存在
	//添加到列表
	buf, find := invertedIdxDB.Get(key)
	ids := make([]uint32, 0)
	if find {
		utils.Decoder(buf, &ids)
		return ids, true
	}

	return nil, false
}

// 处理过滤词文档
func (e *Engine) produceFilterIdx(filteridxChannel chan uint32, filterWord string, wg *sync.WaitGroup) {
	defer wg.Done()
	fileIdxs, ok := e.GetInvertedIdxListByWord(filterWord)
	if !ok {
		// 没有该过滤词相关倒排索引 直接返回
		return
	}
	for _, fileIdx := range fileIdxs {
		filteridxChannel <- fileIdx
	}
}
