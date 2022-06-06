package routers

import (
	"GoSearchEngine/handlers/collection"
	"GoSearchEngine/handlers/dir"
	"GoSearchEngine/handlers/user"
	"GoSearchEngine/middleware"
	"GoSearchEngine/tools"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHTMLMain(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", nil)
}
func GetHTMLIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
func GetHTMLRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}
func GetHTMLLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
func GetHTMLCollectFiles(c *gin.Context) {
	c.HTML(http.StatusOK, "collectFiles.html", nil)
}
func SetupRouter() *gin.Engine {
	//r.LoadHTMLFiles("index.html")
	//r.GET("",index)
	r := gin.Default()
	r.Use(gin.Logger())
	r.LoadHTMLGlob("./web/*")
	r.GET("/main", GetHTMLMain)
	r.GET("/index", GetHTMLIndex)
	r.GET("/register", GetHTMLRegister)
	r.GET("/login", GetHTMLLogin)
	r.GET("/collectFiles", GetHTMLCollectFiles)
	v1Group := r.Group("v1", tools.Cors())
	{
		// 成员管理
		//g.GET("", getUser)                                     // 根据userID获取用户信息
		//g.GET("/user/list", getUserList)                       // 批量获取用户信息（所有用户）
		//g.POST("/user/update", cookieMiddleware(), updateUser) // 更新用户信息

		// 登录
		v1Group.POST("/user/signup", user.Signup)                                // 注册
		v1Group.POST("/user/login", user.Login)                                  // 登录
		v1Group.POST("/user/logout", middleware.CookieMiddleware(), user.Logout) // 登出
		v1Group.GET("/user/whoami", middleware.CookieMiddleware(), user.Whoami)  // 查询当前用户
		//
		//收藏夹
		v1Group.POST("/dir/new", middleware.CookieMiddleware(), dir.NewDir)               // 新建收藏夹
		v1Group.POST("/dir/delete", middleware.CookieMiddleware(), dir.DeleteDir)         // 删除收藏夹
		v1Group.POST("/dir/name", middleware.CookieMiddleware(), dir.UpdateName)          // 更新收藏夹名字
		v1Group.GET("/dir/get", middleware.CookieMiddleware(), dir.GetDir)                //获取用户的收藏夹
		v1Group.POST("/dir/collection", middleware.CookieMiddleware(), dir.GetCollection) // 读取某个收藏夹的内容

		//
		//收藏
		v1Group.POST("/collection/collect", middleware.CookieMiddleware(), collection.Collect)      //收藏
		v1Group.POST("/collection/cancel", middleware.CookieMiddleware(), collection.CancelCollect) //取消收藏

	}
	return r
}
