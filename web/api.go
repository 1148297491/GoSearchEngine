package web

import (
	"gofound/searcher"
	"gofound/searcher/model"
	"gofound/searcher/system"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
)

type Api struct {
	Engine   *searcher.Engine
	Callback func() map[string]interface{}
}

func (a *Api) query(c *gin.Context) {

	var request = &model.SearchRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}

	//调用搜索
	r := a.Engine.MultiSearch(request)
	c.JSON(200, Success(r))
}

func (a *Api) gc(c *gin.Context) {
	runtime.GC()

	c.JSON(200, Success(nil))
}

// status 获取服务器状态
func (a *Api) status(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	s := a.Callback()

	r := gin.H{
		"memory": system.GetMemStat(),
		"cpu":    system.GetCPUStatus(),
		"disk":   system.GetDiskStat(),
		"system": s,
	}
	// 获取服务器状态
	c.JSON(200, Success(r))
}

func (a *Api) addIndex(c *gin.Context) {
	document := &model.IndexDoc{}
	err := c.BindJSON(&document)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}

	go a.Engine.IndexDocument(document) // 引擎管道的生产者

	c.JSON(200, Success(nil))
}

func (a *Api) batchAddIndex(c *gin.Context) {
	documents := make([]model.IndexDoc, 0)
	err := c.BindJSON(&documents)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}

	engine := a.Engine
	for _, doc := range documents {
		go engine.IndexDocument(&doc)
	}

	c.JSON(200, Success(nil))
}

func (a *Api) wordCut(c *gin.Context) {
	q := c.Query("q")
	r := a.Engine.Tokenizer.Cut(q)
	c.JSON(200, Success(r))

}

func welcome(c *gin.Context) {
	c.JSON(200, Success("Welcome to GoFound"))
}

func (a *Api) removeIndex(c *gin.Context) {
	removeIndexModel := &model.RemoveIndexModel{}
	err := c.BindJSON(&removeIndexModel)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}
	engine := a.Engine

	err = engine.RemoveIndex(removeIndexModel.Id)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}
	c.JSON(200, Success(nil))
}

func (a *Api) dbs(ctx *gin.Context) {
	ctx.JSON(200, Success(a.Engine))
}

func (a *Api) restart(c *gin.Context) {

	os.Exit(0)
}

func (a *Api) Register(router *gin.Engine, handlers ...gin.HandlerFunc) {

	group := router.Group("/api", handlers...)

	group.GET("/", welcome)

	group.POST("/query", a.query)

	group.GET("/status", a.status)

	group.GET("/gc", a.gc)

	group.GET("/db/list", a.dbs)

	group.GET("/word/cut", a.wordCut)

	group.POST("/index", a.addIndex)

	group.POST("/index/batch", a.batchAddIndex)

	group.POST("/remove", a.removeIndex)

	group.GET("/restart", a.restart)

}
