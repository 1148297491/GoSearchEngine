package sorts

import (
	"GoSearchEngine/searcher/model"
	"GoSearchEngine/searcher/utils"
	"sort"
	"sync"
)

const (
	DESC = "desc"
)

type ScoreSlice []model.SliceItem

func (x ScoreSlice) Len() int {
	return len(x)
}
func (x ScoreSlice) Less(i, j int) bool {
	return x[i].Score < x[j].Score
}
func (x ScoreSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

type SortSlice []uint32

func (x SortSlice) Len() int {
	return len(x)
}
func (x SortSlice) Less(i, j int) bool {
	return x[i] < x[j]
}
func (x SortSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]

}

type Uint32Slice []uint32

func (x Uint32Slice) Len() int           { return len(x) }
func (x Uint32Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Uint32Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type FastSort struct {
	sync.Mutex

	IsDebug bool

	FilterIdxSet utils.Set

	DataChannel chan uint32

	data []model.SliceItem

	temps []uint32

	count int //总数

	Order string //排序方式
}

func (f *FastSort) Add(ids *[]uint32) {
	for _, id := range *ids {
		if f.FilterIdxSet.IsExist(id) {
			continue
		}
		f.DataChannel <- id
	}
}

func (f *FastSort) AppendData(wg *sync.WaitGroup) {
	defer wg.Done() //等待管道关闭
	for id := range f.DataChannel {
		f.temps = append(f.temps, id)
	}
}

// Count 获取数量
func (f *FastSort) Count() int {
	return f.count
}

// Process 处理数据
func (f *FastSort) Process() {
	// 文档出现的次数越多 说明包含的关键词越多
	fileIDMap := make(map[uint32]int)

	for _, fileId := range f.temps {
		fileIDMap[fileId] += 1
	}

	for k, v := range fileIDMap {
		f.data = append(f.data, model.SliceItem{
			Id:    k,
			Score: v,
		})
	}

	f.count = len(f.data)
	//对分数进行排序
	sort.Sort(sort.Reverse(ScoreSlice(f.data)))
}

func (f *FastSort) GetAll(result *[]model.SliceItem, start int, end int) {
	//获得start到end的data内容
	*result = f.data[start:end]
}
