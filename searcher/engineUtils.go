package searcher

/*
	搜索相关的功能在这个文档，删除等
*/

import (
	"fmt"
	"gofound/searcher/storage"
	"gofound/searcher/utils"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/syndtr/goleveldb/leveldb/errors"
)

// DeleteDuplicatedWordAndCut 对过滤单词进行进一步分词并进行去重 保证和搜索语句分词结果一致
func (e *Engine) DeleteDuplicatedWordAndCut(words []string) []string {
	wordSet := make(utils.Set)
	var res []string = make([]string, 0, len(words)/2)
	for _, word := range words {
		cutWords := e.Tokenizer.Cut(word, true) // 采用精准分词模式
		for _, cutWord := range cutWords {
			if wordSet.IsExist(cutWord) || word == "" {
				continue
			}

			wordSet.Add(word)
			res = append(res, word)
		}
	}

	return res
}

// 去掉过滤词第一步 首先对字符串包含的
func (e *Engine) deleteFilterWordsFromQuery(filterWords []string, query string) string {
	for idx, filterWord := range filterWords {
		if strings.Contains(query, filterWord) {
			query = strings.ReplaceAll(query, filterWord, "")
			filterWords[idx] = ""
		}
	}
	return query
}

// GetIndexCount 获取索引数量
func (e *Engine) GetIndexCount() int64 {
	var size int64
	for i := 0; i < e.Shard; i++ {
		size += e.invertedIndexStorages[i].Count()
	}
	return size
}

// GetDocumentCount 获取文档数量
func (e *Engine) GetDocumentCount() int64 {
	var count int64
	for i := 0; i < e.Shard; i++ {
		count += e.docStorages[i].Count()
	}
	return count
}

// GetDocById 通过id获取文档
func (e *Engine) GetDocById(id uint32) []byte {
	shard := e.getShard(id)
	key := utils.Uint32ToBytes(id)
	buf, found := e.docStorages[shard].Get(key)
	if found {
		return buf
	}

	return nil
}

// RemoveIndex 根据ID移除索引
func (e *Engine) RemoveIndex(id uint32) error {
	//移除
	e.Lock()
	defer e.Unlock()

	shard := e.getShard(id)
	key := utils.Uint32ToBytes(id)

	//关键字和Id映射
	//invertedIndexStorages []*storage.LeveldbStorage
	//ID和key映射，用于计算相关度，一个id 对应多个key
	ik := e.positiveIndexStorages[shard]
	keysValue, found := ik.Get(key)
	if !found {
		return errors.New(fmt.Sprintf("没有找到id=%d", id))
	}

	keys := make([]string, 0)
	utils.Decoder(keysValue, &keys)

	//符合条件的key，要移除id
	for _, word := range keys {
		e.removeIdInWordIndex(id, word)
	}

	//删除id映射
	err := ik.Delete(key)
	if err != nil {
		return errors.New(err.Error())
	}

	//文档仓
	err = e.docStorages[shard].Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) Close() {
	e.Lock()
	defer e.Unlock()

	for i := 0; i < e.Shard; i++ {
		e.invertedIndexStorages[i].Close()
		e.positiveIndexStorages[i].Close()
	}
}

// Drop 删除 == 删除本地db库，注意需要加管理员权限
func (e *Engine) Drop() error {
	e.Lock()
	defer e.Unlock()
	//删除文件
	dir, err := ioutil.ReadDir(e.IndexPath)
	if err != nil {
		return err
	}
	for _, d := range dir {
		err := os.RemoveAll(path.Join([]string{d.Name()}...))
		if err != nil {
			return err
		}
		os.Remove(e.IndexPath)
	}

	//清空内存
	for i := 0; i < e.Shard; i++ {
		e.docStorages = make([]*storage.LeveldbStorage, 0)
		e.invertedIndexStorages = make([]*storage.LeveldbStorage, 0)
		e.positiveIndexStorages = make([]*storage.LeveldbStorage, 0)
	}

	return nil
}
