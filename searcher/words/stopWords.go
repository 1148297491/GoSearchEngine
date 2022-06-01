package words

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var stopWordsMap map[string]bool

func NewStopWords(filePath string) {
	stopWordsMap = make(map[string]bool, 0)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	var size = stat.Size()
	fmt.Println("stopWordsDict size =", size)
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.ReplaceAll(line, "\n", "")
		stopWordsMap[line] = true
		if err != nil {
			if err == io.EOF {
				fmt.Println("read end!")
			} else {
				panic(err)
			}
			break
		}
	}

}

func GetStopWordsMap() map[string]bool {
	//sfmt.Println(stopWordsMap["是"])
	return stopWordsMap
}
