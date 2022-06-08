package words

import (
	"GoSearchEngine/searcher/utils"
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var stopWordsSet utils.Set

func init() {
	log.Printf("###### 开始初始化去停用词表 ######")
	_, fileStr, _, _ := runtime.Caller(0)

	currPtah := filepath.Dir(fileStr)
	var stopWordsFilePath = filepath.Join(currPtah, "/data/stopWordsDict.txt")

	file, err := os.Open(stopWordsFilePath)
	if err != nil {
		log.Printf("[init] 打开停用词典失败, err=%+v", err)
		panic(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Printf("[init] 获取文件状态失败, err=%+v", err)
		panic(err)
	}

	var size = stat.Size()
	log.Printf("停用词典大小为: %+v", size)
	stopWordsSet = make(utils.Set, size)

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("[init] 去停用词表独取完成")
			} else {
				log.Printf("[init] 独取去停用词表发生错误, err=%+v", err)
				panic(err)
			}
			break
		}

		line = strings.ReplaceAll(line, "\n", "")
		stopWordsSet.Add(line)
	}

	log.Printf("###### 初始化去停用词表完成 ######")
}

func GetStopWordsSet() utils.Set {
	return stopWordsSet
}
