# SuperGoFound

## 搜索请求结构体如下
```
type SearchRequest struct {
	Query       string   `json:"query,omitempty"`      // 搜索关键词
	Order       string   `json:"order,omitempty"`      // 排序类型
	Page        int      `json:"page,omitempty"`       // 页码
	Limit       int      `json:"limit,omitempty"`      // 每页大小，最大1000，超过报错
	FilterWords []string `json:"filterword,omitempty"` //用户屏蔽词列表
}
```
## 搜索结果结构体如下
```
type SearchResult struct {
	Time          float64       `json:"time,omitempty"`          //查询用时
	Total         int           `json:"total"`                   //总数
	PageCount     int           `json:"pageCount"`               //总页数
	Page          int           `json:"page,omitempty"`          //页码
	Limit         int           `json:"limit,omitempty"`         //页大小
	Documents     []ResponseDoc `json:"documents,omitempty"`     //文档
	RelatedSearch []string      `json:"relatedSearch,omitempty"` // 相关搜素
}
```
## 词典存储位置
`/searcher/words/data/dictionary.txt`

## 运行流程
```
go run ./main.go
```

## POST请求示例
```
curl -H "Content-Type: application/json" -X POST -d '{"query":"番茄", "order":"", "page":1, "limit":1000, "filterwords":"鸡蛋"}' "127.0.0.1:9633/search"
```
