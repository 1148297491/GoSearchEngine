package web

import (
	"gofound/searcher"
	"gofound/searcher/model"

	"github.com/gin-gonic/gin"
)

type WebSearch struct {
	SearchEngine *searcher.Engine
}

func (e *WebSearch) SearchProcessFun(c *gin.Context) {
	var request = &model.SearchRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}

	//调用搜索
	r := e.SearchEngine.MultiSearch(request)
	c.JSON(200, Success(r))
}
