package initwukong

import (
	"GoSearchEngine/searcher"
	"GoSearchEngine/searcher/model"
	"bufio"
	"bytes"
	"encoding/csv"

	"io"
	"log"
	"os"
	"sync"
)

/*
	并发读取 wukong数据集 提升初始化速度
*/
const MB = 1024 * 1024
const GB = 1024 * MB

// 设立文件缓冲区对象池
var readBufPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 10*MB)
		return &buf
	},
}

type read struct {
	idChannels             []chan uint32
	limitProcessNumChannel chan struct{}
	isProduceId            bool
	routineNum             int
}

func NewReadInstance(goRoutineNum int) *read {
	readInstance := &read{
		idChannels:             make([]chan uint32, goRoutineNum),
		limitProcessNumChannel: make(chan struct{}, goRoutineNum),
		isProduceId:            true,
		routineNum:             goRoutineNum,
	}

	for i := 0; i < goRoutineNum; i++ {
		readInstance.idChannels[i] = make(chan uint32, 10)
	}

	for startId := 0; startId < goRoutineNum; startId++ {
		go readInstance.produceId(startId, goRoutineNum)
	}

	return readInstance
}

func (r *read) produceId(startId, goRoutineNum int) {
	defer close(r.idChannels[startId])

	var id uint32 = uint32(startId)

	for ; ; id += uint32(goRoutineNum) {
		if !r.isProduceId {
			break
		}

		r.idChannels[startId] <- id

	}
	log.Printf("不再生产id")
}

func (r *read) MultiRead(filePath string, engine *searcher.Engine, fileSize int) {
	wukong, err := os.Open(filePath)

	if err != nil {
		log.Printf("打开wukong数据集出现问题, err=%+v", err)
		panic(err)
	}
	defer wukong.Close()
	defer close(r.limitProcessNumChannel)

	// 按照缓冲区的大小分块读取文件 管道控制并发大小
	fileBlockNum := 0
	limitChannelCap := cap(r.limitProcessNumChannel)
	wg := &sync.WaitGroup{}
	wukongReader := bufio.NewReader(wukong)
	wukongReader = bufio.NewReaderSize(wukongReader, fileSize*MB)

	wukongReader.ReadLine() // 跳过第一行

	// 开始分块读取文件
	for {
		readBuf := readBufPool.Get().(*[]byte)

		byteSize, err := wukongReader.Read(*readBuf)

		if err != nil {
			if err == io.EOF {
				log.Printf("分块读取文件完毕!")
			} else {
				log.Printf("分块读取文件出现错误, err=%+v", err)
			}

			break
		}

		//log.Printf("实际读取文件块大小{%+v}", byteSize)

		buf := (*readBuf)[:byteSize]

		// 考虑读取的byte不一定是一行的结束
		remainBytesUtilNewLine, err := wukongReader.ReadBytes('\n')

		if err != nil && err != io.EOF {
			log.Printf("读取文件遇见错误err=%+v", err)
			break
		}

		buf = append(buf, remainBytesUtilNewLine...)

		// 限制一次最多处理 goRoutineNum 大小的内存块
		r.limitProcessNumChannel <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				<-r.limitProcessNumChannel
			}()

			r.processReadBuf(buf, engine, fileBlockNum%limitChannelCap)
		}()
		fileBlockNum += 1
	}

	// 等待所有的协程结束 关闭id生产管道
	wg.Wait()
	r.isProduceId = false
}

func (r *read) processReadBuf(readBuf []byte, engine *searcher.Engine, idChannelIdx int) {
	// 对读取的byte文件块进行处理
	csvFile := csv.NewReader(bytes.NewReader(readBuf))

	// 使用结束之后放回缓冲区词中
	readBufPool.Put(&readBuf)
	var id = uint32(0)

	for {
		urlAndText, err := csvFile.Read()
		if err == io.EOF {
			// 读到文件结尾退出
			break
		} else if err != nil {
			log.Printf("读取单条csv错误: %+v", err)
			continue
		}

		// 获取id
		id = <-r.idChannels[idChannelIdx]

		// 构建索引实体
		indexInstance := &model.IndexDoc{
			Id:       id,
			Text:     urlAndText[1],
			Document: make(map[string]interface{}, 3),
		}

		indexInstance.Document["id"] = id
		indexInstance.Document["url"] = urlAndText[0]
		indexInstance.Document["text"] = urlAndText[1]

		engine.IndexDocument(indexInstance)
	}
}

func (r *read) splitCsvLine(line []byte) []string {
	// url中没有逗号 直接按照第一个逗号位置进行切分
	if len(line) == 0 {
		return nil
	}

	var csvRes = make([]string, 0, 2)
	if line[0] != '"' {
		pos := bytes.IndexByte(line, ',') // 查找第一个逗号位置
		csvRes = append(csvRes, string(line[:pos]))
		csvRes = append(csvRes, string(line[pos+1:]))
		return csvRes
	}

	// 处理如果url中含有逗号
	pos := bytes.IndexByte(line[1:], '"')
	pos = pos + 2
	csvRes = append(csvRes, string(line[:pos]))
	csvRes = append(csvRes, string(line[pos+1:]))

	return csvRes
}
