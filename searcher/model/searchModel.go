package model

// IndexDoc 索引实体
type IndexDoc struct {
	Id       uint32                 `json:"id,omitempty"`
	Text     string                 `json:"text,omitempty"`
	Document map[string]interface{} `json:"document,omitempty"`
}

// StorageIndexDoc 文档对象
type StorageIndexDoc struct {
	*IndexDoc
	Keys []string `json:"keys,omitempty"`
}

//搜索结果返回结构体
type ResponseDoc struct {
	Score         int    `json:"score,omitempty"` //得分
	HighlightText string `json:"highlightText,omitempty"`
	Url           string `json:"url,omitempty"`
	OriginalText  string `json:"originalText,omitempty"`
}

//删除索引文章id
type RemoveIndexModel struct {
	Id uint32 `json:"id,omitempty"`
}

//搜索结果id及分数（结果=结构体）
type SliceItem struct {
	Id    uint32
	Score int
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query       string   `json:"query,omitempty"`      // 搜索关键词
	Order       string   `json:"order,omitempty"`      // 排序类型
	Page        int      `json:"page,omitempty"`       // 页码
	Limit       int      `json:"limit,omitempty"`      // 每页大小，最大1000，超过报错
	FilterWords []string `json:"filterword,omitempty"` //用户屏蔽词列表
}

func (s *SearchRequest) GetAndSetDefault() *SearchRequest {

	if s.Limit == 0 {
		s.Limit = 100
	}
	if s.Page == 0 {
		s.Page = 1
	}

	if s.Order == "" {
		s.Order = "desc"
	}

	return s
}

// SearchResult 搜索响应
type SearchResult struct {
	Time          float64       `json:"time,omitempty"`          //查询用时
	Total         int           `json:"total"`                   //总数
	PageCount     int           `json:"pageCount"`               //总页数
	Page          int           `json:"page,omitempty"`          //页码
	Limit         int           `json:"limit,omitempty"`         //页大小
	Documents     []ResponseDoc `json:"documents,omitempty"`     //文档
	RelatedSearch []string      `json:"relatedSearch,omitempty"` // 相关搜素
}
