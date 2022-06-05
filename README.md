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
修改config.yml并创建数据库
go run ./main.go
```
## 前端访问
```
http://127.0.0.1:9633/main
```

## POST请求示例
```
curl -H "Content-Type: application/json" -X POST -d '{"query":"番茄", "order":"", "page":1, "limit":1000, "filterword":["鸡蛋"]}' "127.0.0.1:9633/search"
```


## 以图搜图部分：
下载这里面的数据文件
链接：https://pan.baidu.com/s/1e_8G4ou_9Dv-Jvshlxddhw 
提取码：1612

安装python3.8，pytorch，faiss1.7.2  (pip --default-time=1000 install -i https://pypi.tuna.tsinghua.edu.cn/simple faiss-cpu==1.7.2)，之后修改findGraph/flask_app.py中app.config['UPLOAD_FOLDER']
和findGraph/predict.py中数据文件的路径
```
index = faiss.read_index('index_file.index')
with open('wukong_urls.txt', 'r', encoding='utf-8') as f:
    urlsTotal = [row.strip() for row in f.readlines()]
```
之后运行flask_app.py即可
